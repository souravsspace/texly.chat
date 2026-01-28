package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
)

/*
* EmbeddingService handles generating vector embeddings using OpenAI API
 */
type EmbeddingService struct {
	apiKey     string
	model      string
	dimensions int
	httpClient *http.Client
}

/*
* NewEmbeddingService creates a new embedding service instance
 */
func NewEmbeddingService(apiKey, model string, dimensions int) *EmbeddingService {
	return &EmbeddingService{
		apiKey:     apiKey,
		model:      model,
		dimensions: dimensions,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

/*
* OpenAI API request/response structures
 */
type embeddingRequest struct {
	Input      interface{} `json:"input"`
	Model      string      `json:"model"`
	Dimensions int         `json:"dimensions,omitempty"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

/*
* GenerateEmbedding generates a single embedding vector for the given text
 */
func (s *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := s.GenerateEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return embeddings[0], nil
}

/*
* GenerateEmbeddings generates embeddings for multiple texts in a batch
* OpenAI supports up to 2048 inputs per request
 */
func (s *EmbeddingService) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	if len(texts) > 2048 {
		return nil, fmt.Errorf("too many texts: %d (max 2048)", len(texts))
	}

	// Prepare request
	reqBody := embeddingRequest{
		Input:      texts,
		Model:      s.model,
		Dimensions: s.dimensions,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = s.httpClient.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("failed to execute request after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
			continue
		}
		break
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	// Parse successful response
	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract embeddings
	embeddings := make([][]float32, len(texts))
	for _, data := range embResp.Data {
		if data.Index < 0 || data.Index >= len(embeddings) {
			return nil, fmt.Errorf("invalid embedding index: %d", data.Index)
		}
		embeddings[data.Index] = data.Embedding
	}

	return embeddings, nil
}

/*
* EmbedChunks is a convenience method to generate embeddings for document chunks
 */
func (s *EmbeddingService) EmbedChunks(ctx context.Context, chunks []models.DocumentChunk) ([][]float32, error) {
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}
	return s.GenerateEmbeddings(ctx, texts)
}
