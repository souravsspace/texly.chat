package chunker

import (
	"strings"
	"unicode"
)

/*
* ChunkText splits text into chunks of approximately maxTokens size
* This uses a simple word-count approximation (tokens ≈ words * 1.3)
 */
func ChunkText(content string, maxTokens int) []string {
	// Approximate: 1 token ≈ 0.75 words, so maxWords ≈ maxTokens * 0.75
	maxWords := int(float64(maxTokens) * 0.75)
	
	// Split into paragraphs first
	paragraphs := strings.Split(content, "\n")
	
	var chunks []string
	var currentChunk strings.Builder
	currentWordCount := 0
	
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		
		paraWords := countWords(para)
		
		// If adding this paragraph would exceed limit, save current chunk and start new one
		if currentWordCount > 0 && currentWordCount+paraWords > maxWords {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
			currentWordCount = 0
		}
		
		// If a single paragraph is too large, split it by sentences
		if paraWords > maxWords {
			sentences := splitSentences(para)
			for _, sentence := range sentences {
				sentenceWords := countWords(sentence)
				
				if currentWordCount > 0 && currentWordCount+sentenceWords > maxWords {
					chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
					currentChunk.Reset()
					currentWordCount = 0
				}
				
				if currentChunk.Len() > 0 {
					currentChunk.WriteString(" ")
				}
				currentChunk.WriteString(sentence)
				currentWordCount += sentenceWords
			}
		} else {
			// Add paragraph to current chunk
			if currentChunk.Len() > 0 {
				currentChunk.WriteString("\n\n")
			}
			currentChunk.WriteString(para)
			currentWordCount += paraWords
		}
	}
	
	// Add any remaining content
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}
	
	return chunks
}

/*
* countWords counts the number of words in a text
 */
func countWords(text string) int {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return len(words)
}

/*
* splitSentences splits text into sentences (simple implementation)
 */
func splitSentences(text string) []string {
	// Simple sentence splitting by common punctuation
	text = strings.ReplaceAll(text, ". ", ".|")
	text = strings.ReplaceAll(text, "! ", "!|")
	text = strings.ReplaceAll(text, "? ", "?|")
	
	sentences := strings.Split(text, "|")
	var result []string
	
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}
	
	return result
}
