package extractor

import (
	"fmt"
	"io"
	"strings"
)

/*
* TextReader handles plain text and markdown file reading
 */
type TextReader struct{}

/*
* NewTextReader creates a new text reader instance
 */
func NewTextReader() *TextReader {
	return &TextReader{}
}

/*
* ReadTextFile reads plain text or markdown file content
 */
func (r *TextReader) ReadTextFile(reader io.Reader) (string, error) {
	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	text := strings.TrimSpace(string(content))
	if len(text) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	return text, nil
}
