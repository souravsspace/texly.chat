package extractor

import (
	"strings"
	"testing"
)

func TestTextReader_ReadTextFile(t *testing.T) {
	reader := NewTextReader()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid text content",
			content: "This is a test document.\nWith multiple lines.",
			wantErr: false,
		},
		{
			name:    "valid markdown content",
			content: "# Heading\n\nThis is markdown content.\n\n- List item 1\n- List item 2",
			wantErr: false,
		},
		{
			name:    "empty file",
			content: "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			content: "   \n\n  \t  ",
			wantErr: true,
		},
		{
			name:    "single line",
			content: "Single line of text",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.content)
			text, err := reader.ReadTextFile(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTextFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(text) == 0 {
					t.Error("ReadTextFile() returned empty text for valid content")
				}
				if !strings.Contains(text, strings.TrimSpace(tt.content)) {
					t.Error("ReadTextFile() did not preserve content")
				}
			}
		})
	}
}

func TestTextReader_ReadTextFile_LargeContent(t *testing.T) {
	reader := NewTextReader()

	// Create large content (1MB)
	largeContent := strings.Repeat("This is a line of text.\n", 50000)
	r := strings.NewReader(largeContent)

	text, err := reader.ReadTextFile(r)
	if err != nil {
		t.Errorf("ReadTextFile() failed on large content: %v", err)
	}

	if len(text) == 0 {
		t.Error("ReadTextFile() returned empty text for large content")
	}
}
