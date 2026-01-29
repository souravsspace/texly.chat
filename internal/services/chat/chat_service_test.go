package chat

import (
	"context"
	"testing"

	"github.com/souravsspace/texly.chat/internal/services/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * Test ChatService initialization
 */
func TestNewChatService(t *testing.T) {
	service := NewChatService(
		nil, // embedding service
		nil, // search service
		"gpt-4o-mini",
		0.7,
		5,
		"test-api-key",
	)

	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
	assert.Equal(t, "gpt-4o-mini", string(service.chatModel))
	assert.Equal(t, 0.7, service.temperature)
	assert.Equal(t, 5, service.maxContextChunks)
}

/*
 * Test buildMessages with no context
 */
func TestBuildMessages_NoContext(t *testing.T) {
	service := NewChatService(nil, nil, "gpt-4o-mini", 0.7, 5, "test-key")

	messages := service.buildMessages(
		"You are a helpful assistant",
		[]vector.SearchResult{},
		"Hello, how are you?",
	)

	// Should have system + user message
	assert.Len(t, messages, 2)
}

/*
 * Test buildMessages with context
 */
func TestBuildMessages_WithContext(t *testing.T) {
	service := NewChatService(nil, nil, "gpt-4o-mini", 0.7, 5, "test-key")

	contextChunks := []vector.SearchResult{
		{
			ChunkID:    "chunk1",
			Content:    "This is relevant context",
			URL:        "https://example.com/doc1",
			Distance:   0.1,
			ChunkIndex: 0,
		},
		{
			ChunkID:    "chunk2",
			Content:    "More relevant information",
			URL:        "https://example.com/doc2",
			Distance:   0.2,
			ChunkIndex: 1,
		},
	}

	messages := service.buildMessages(
		"You are a helpful assistant",
		contextChunks,
		"Tell me about the docs",
	)

	// Should have system + context + user message
	assert.Len(t, messages, 3)
}

/*
 * Test buildMessages without system prompt
 */
func TestBuildMessages_NoSystemPrompt(t *testing.T) {
	service := NewChatService(nil, nil, "gpt-4o-mini", 0.7, 5, "test-key")

	messages := service.buildMessages(
		"", // No system prompt
		[]vector.SearchResult{},
		"Hello!",
	)

	// Should only have user message
	assert.Len(t, messages, 1)
}

/*
 * Test StreamChat with context cancellation
 */
func TestStreamChat_ContextCancellation(t *testing.T) {
	// Skip this test if OpenAI API key is not set
	// This is an integration test that requires actual API access
	t.Skip("Skipping integration test - requires OpenAI API key")

	service := NewChatService(nil, nil, "gpt-4o-mini", 0.7, 5, "test-key")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	tokenChan, errChan := service.StreamChat(
		ctx,
		"bot-123",
		"You are helpful",
		"Hello",
	)

	// Should receive error due to cancelled context
	select {
	case <-tokenChan:
		t.Fatal("Expected no tokens from cancelled context")
	case err := <-errChan:
		require.Error(t, err)
	}
}

/*
 * Test buildMessages formatting with multiple context chunks
 */
func TestBuildMessages_ContextFormatting(t *testing.T) {
	service := NewChatService(nil, nil, "gpt-4o-mini", 0.7, 5, "test-key")

	contextChunks := []vector.SearchResult{
		{
			Content: "First chunk content",
			URL:     "https://example.com/page1",
		},
		{
			Content: "Second chunk content",
			URL:     "https://example.com/page2",
		},
		{
			Content: "Third chunk content",
			URL:     "https://example.com/page3",
		},
	}

	messages := service.buildMessages(
		"System prompt",
		contextChunks,
		"User question",
	)

	// Verify structure: system + context + user = 3 messages
	assert.Len(t, messages, 3)
}

/*
 * Test ChatService with different temperature values
 */
func TestChatService_TemperatureConfiguration(t *testing.T) {
	testCases := []struct {
		name        string
		temperature float64
	}{
		{"Low temperature", 0.0},
		{"Medium temperature", 0.7},
		{"High temperature", 1.5},
		{"Max temperature", 2.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewChatService(
				nil,
				nil,
				"gpt-4o-mini",
				tc.temperature,
				5,
				"test-key",
			)

			assert.Equal(t, tc.temperature, service.temperature)
		})
	}
}

/*
 * Test ChatService with different max context chunks
 */
func TestChatService_MaxContextChunks(t *testing.T) {
	testCases := []struct {
		name             string
		maxContextChunks int
	}{
		{"No context", 0},
		{"Single chunk", 1},
		{"Default chunks", 5},
		{"Many chunks", 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewChatService(
				nil,
				nil,
				"gpt-4o-mini",
				0.7,
				tc.maxContextChunks,
				"test-key",
			)

			assert.Equal(t, tc.maxContextChunks, service.maxContextChunks)
		})
	}
}
