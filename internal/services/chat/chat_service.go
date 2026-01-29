package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/vector"
)

/*
 * ChatService orchestrates RAG-powered chat with streaming responses
 */
type ChatService struct {
	embeddingService *embedding.EmbeddingService
	searchService    *vector.SearchService
	chatModel        openai.ChatModel
	temperature      float64
	maxContextChunks int
	client           openai.Client
}

/*
 * NewChatService creates a new chat service instance
 */
func NewChatService(
	embeddingService *embedding.EmbeddingService,
	searchService *vector.SearchService,
	chatModel string,
	temperature float64,
	maxContextChunks int,
	apiKey string,
) *ChatService {
	return &ChatService{
		embeddingService: embeddingService,
		searchService:    searchService,
		chatModel:        openai.ChatModel(chatModel),
		temperature:      temperature,
		maxContextChunks: maxContextChunks,
		client:           openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

/*
 * StreamChat performs RAG and streams LLM response via channels
 * Returns a token channel and error channel
 */
func (s *ChatService) StreamChat(
	ctx context.Context,
	botID string,
	systemPrompt string,
	userMessage string,
) (<-chan string, <-chan error) {
	tokenChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errChan)

		// Step 1: Perform RAG - retrieve relevant context
		contextChunks, err := s.searchService.SearchSimilar(ctx, userMessage, botID, s.maxContextChunks)
		if err != nil {
			errChan <- fmt.Errorf("failed to search context: %w", err)
			return
		}

		// Step 2: Build messages with context
		messages := s.buildMessages(systemPrompt, contextChunks, userMessage)

		// Step 3: Stream from OpenAI using official SDK
		if err := s.streamFromOpenAI(ctx, messages, tokenChan); err != nil {
			errChan <- err
			return
		}
	}()

	return tokenChan, errChan
}

/*
 * buildMessages constructs the complete message array with system, context, and user message
 */
func (s *ChatService) buildMessages(
	systemPrompt string,
	contextChunks []vector.SearchResult,
	userMessage string,
) []openai.ChatCompletionMessageParamUnion {
	messages := []openai.ChatCompletionMessageParamUnion{}

	// System message
	if systemPrompt != "" {
		messages = append(messages, openai.SystemMessage(systemPrompt))
	}

	// Add context from RAG if available
	if len(contextChunks) > 0 {
		var contextBuilder strings.Builder
		contextBuilder.WriteString("Here is relevant information from the knowledge base:\n\n")

		for i, chunk := range contextChunks {
			contextBuilder.WriteString(fmt.Sprintf("--- Context %d ---\n", i+1))
			contextBuilder.WriteString(chunk.Content)
			contextBuilder.WriteString(fmt.Sprintf("\nSource: %s\n\n", chunk.URL))
		}

		contextBuilder.WriteString("Please use this information to answer the user's question accurately.")

		messages = append(messages, openai.SystemMessage(contextBuilder.String()))
	}

	// User message
	messages = append(messages, openai.UserMessage(userMessage))

	return messages
}

/*
 * streamFromOpenAI makes streaming request to OpenAI Chat Completion API using official SDK
 */
func (s *ChatService) streamFromOpenAI(
	ctx context.Context,
	messages []openai.ChatCompletionMessageParamUnion,
	tokenChan chan<- string,
) error {
	// Create streaming chat completion
	stream := s.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       s.chatModel,
		Temperature: openai.Float(s.temperature),
	})

	// Process the stream
	for stream.Next() {
		chunk := stream.Current()
		
		// Extract content delta from the first choice
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			select {
			case tokenChan <- chunk.Choices[0].Delta.Content:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	// Check for streaming errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("streaming error: %w", err)
	}

	return nil
}
