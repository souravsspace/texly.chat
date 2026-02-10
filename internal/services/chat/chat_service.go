package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/souravsspace/texly.chat/internal/models"
	messageRepo "github.com/souravsspace/texly.chat/internal/repo/message"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/vector"
	"github.com/tiktoken-go/tokenizer"
)

/*
 * ChatService orchestrates RAG-powered chat with streaming responses
 */
type ChatService struct {
	embeddingService *embedding.EmbeddingService
	searchService    *vector.SearchService
	messageRepo      *messageRepo.MessageRepository
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
	messageRepo *messageRepo.MessageRepository,
	chatModel string,
	temperature float64,
	maxContextChunks int,
	apiKey string,
) *ChatService {
	return &ChatService{
		embeddingService: embeddingService,
		searchService:    searchService,
		messageRepo:      messageRepo,
		chatModel:        openai.ChatModel(chatModel),
		temperature:      temperature,
		maxContextChunks: maxContextChunks,
		client:           openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

/*
 * StreamChat performs RAG and streams LLM response via channels
 * Returns a token channel and error channel
 * Optionally saves messages to database if sessionID and userID are provided
 */
func (s *ChatService) StreamChat(
	ctx context.Context,
	botID string,
	systemPrompt string,
	userMessage string,
	sessionID string,
	userID *string,
) (<-chan string, <-chan error) {
	tokenChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errChan)

		// Step 1: Save user message to database
		if s.messageRepo != nil && sessionID != "" {
			userMsg := &models.Message{
				SessionID:  sessionID,
				BotID:      botID,
				UserID:     userID,
				Role:       "user",
				Content:    userMessage,
				TokenCount: countTokens(userMessage),
			}
			if err := s.messageRepo.Create(ctx, userMsg); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Warning: failed to save user message: %v\n", err)
			}
		}

		// Step 2: Perform RAG - retrieve relevant context
		contextChunks, err := s.searchService.SearchSimilar(ctx, userMessage, botID, s.maxContextChunks)
		if err != nil {
			errChan <- fmt.Errorf("failed to search context: %w", err)
			return
		}

		// Step 3: Build messages with context
		messages := s.buildMessages(systemPrompt, contextChunks, userMessage)

		// Step 4: Stream from OpenAI and collect response
		var fullResponse strings.Builder
		if err := s.streamFromOpenAI(ctx, messages, tokenChan, &fullResponse); err != nil {
			errChan <- err
			return
		}

		// Step 5: Save assistant message to database
		if s.messageRepo != nil && sessionID != "" {
			assistantMsg := &models.Message{
				SessionID:  sessionID,
				BotID:      botID,
				UserID:     userID,
				Role:       "assistant",
				Content:    fullResponse.String(),
				TokenCount: countTokens(fullResponse.String()),
			}
			if err := s.messageRepo.Create(ctx, assistantMsg); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Warning: failed to save assistant message: %v\n", err)
			}
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
 * Also collects the full response in responseBuilder for persistence
 */
func (s *ChatService) streamFromOpenAI(
	ctx context.Context,
	messages []openai.ChatCompletionMessageParamUnion,
	tokenChan chan<- string,
	responseBuilder *strings.Builder,
) error {
	// Create streaming chat completion params
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    s.chatModel,
	}

	// Only set temperature if it's not the default (1.0)
	// Some models (like o1) don't support custom temperature
	if s.temperature != 1.0 {
		params.Temperature = openai.Float(s.temperature)
	}

	stream := s.client.Chat.Completions.NewStreaming(ctx, params)

	// Process the stream
	for stream.Next() {
		chunk := stream.Current()

		// Extract content delta from the first choice
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			token := chunk.Choices[0].Delta.Content

			// Collect full response
			if responseBuilder != nil {
				responseBuilder.WriteString(token)
			}

			// Stream token to channel
			select {
			case tokenChan <- token:
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

/*
 * countTokens estimates token count for a string
 * Uses a simple approximation: 1 token ≈ 4 characters
 */
func countTokens(text string) int {
	// Try to use tiktoken for accurate counting
	codec, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		// Fallback to simple approximation
		return len(text) / 4
	}

	tokens, _, err := codec.Encode(text)
	if err != nil {
		// Fallback to simple approximation
		return len(text) / 4
	}

	return len(tokens)
}
