package chunker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkText_SmallContent(t *testing.T) {
	content := "This is a small piece of content."
	maxTokens := 100

	chunks := ChunkText(content, maxTokens)

	assert.Len(t, chunks, 1)
	assert.Equal(t, content, chunks[0])
}

func TestChunkText_LargeParagraphs(t *testing.T) {
	// Create content with multiple paragraphs
	para1 := strings.Repeat("word ", 50) // ~50 words
	para2 := strings.Repeat("text ", 50) // ~50 words
	para3 := strings.Repeat("data ", 50) // ~50 words
	content := para1 + "\n\n" + para2 + "\n\n" + para3

	// Set maxTokens to allow ~40 words per chunk
	maxTokens := 50

	chunks := ChunkText(content, maxTokens)

	// Should create multiple chunks
	assert.Greater(t, len(chunks), 1)

	// Each chunk should not be empty
	for _, chunk := range chunks {
		assert.NotEmpty(t, chunk)
	}
}

func TestChunkText_LongSentences(t *testing.T) {
	// Create a single paragraph with very long sentences
	sentence1 := strings.Repeat("word ", 100) + ". "
	sentence2 := strings.Repeat("text ", 100) + ". "
	content := sentence1 + sentence2

	maxTokens := 80 // ~60 words

	chunks := ChunkText(content, maxTokens)

	// Should split by sentences
	assert.Greater(t, len(chunks), 1)

	// Verify chunks contain complete sentences (end with period)
	for i, chunk := range chunks {
		if i < len(chunks)-1 {
			// Not the last chunk - should contain periods
			assert.Contains(t, chunk, ".")
		}
	}
}

func TestChunkText_EmptyContent(t *testing.T) {
	content := ""
	maxTokens := 100

	chunks := ChunkText(content, maxTokens)

	assert.Empty(t, chunks)
}

func TestChunkText_WhitespaceOnly(t *testing.T) {
	content := "   \n\n   \n   "
	maxTokens := 100

	chunks := ChunkText(content, maxTokens)

	assert.Empty(t, chunks)
}

func TestChunkText_MultipleNewlines(t *testing.T) {
	content := "Paragraph 1\n\n\n\n\nParagraph 2\n\nParagraph 3"
	maxTokens := 100

	chunks := ChunkText(content, maxTokens)

	// Should be combined into one chunk since it's small
	assert.Len(t, chunks, 1)
	// Should clean up excessive newlines
	assert.Contains(t, chunks[0], "Paragraph 1")
	assert.Contains(t, chunks[0], "Paragraph 2")
	assert.Contains(t, chunks[0], "Paragraph 3")
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"hello world", 2},
		{"one", 1},
		{"", 0},
		{"  spaces   between   words  ", 3},
		{"multiple\nline\ntext", 3},
		{"punctuation, and! symbols?", 3},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			count := countWords(tt.text)
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestSplitSentences(t *testing.T) {
	text := "First sentence. Second sentence! Third sentence? Fourth sentence"
	sentences := splitSentences(text)

	assert.GreaterOrEqual(t, len(sentences), 3)
	assert.Contains(t, strings.Join(sentences, " "), "First sentence")
	assert.Contains(t, strings.Join(sentences, " "), "Second sentence")
	assert.Contains(t, strings.Join(sentences, " "), "Third sentence")
}

func TestChunkText_RealWorldExample(t *testing.T) {
	content := `# Introduction

This is the introduction paragraph with some content. It explains the main topic.

## Section 1

This is section 1 with detailed information. It has multiple sentences. Each sentence adds more context.

## Section 2

This is section 2 with even more content. The content continues here with additional details. More information follows.`

	maxTokens := 50 // Small chunks

	chunks := ChunkText(content, maxTokens)

	// Should create multiple chunks
	assert.Greater(t, len(chunks), 1)

	// Verify all content is included
	fullContent := strings.Join(chunks, " ")
	assert.Contains(t, fullContent, "Introduction")
	assert.Contains(t, fullContent, "Section 1")
	assert.Contains(t, fullContent, "Section 2")
}
