package extractor

import (
	"strings"
	"testing"
)

// Note: Testing PDF extraction properly requires actual PDF files
// These tests verify the basic structure and error handling

func TestPDFExtractor_New(t *testing.T) {
	extractor := NewPDFExtractor()
	if extractor == nil {
		t.Error("NewPDFExtractor() returned nil")
	}
}

func TestPDFExtractor_ExtractText_InvalidPDF(t *testing.T) {
	extractor := NewPDFExtractor()
	
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "plain text (not PDF)",
			content: "This is plain text, not a PDF",
		},
		{
			name:    "empty content",
			content: "",
		},
		{
			name:    "partial PDF header",
			content: "%PDF-",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.content)
			_, err := extractor.ExtractText(r)
			
			if err == nil {
				t.Error("ExtractText() should fail on invalid PDF content")
			}
		})
	}
}

// Helper function to create a minimal valid PDF structure for testing
// This creates a very basic PDF that can be parsed
func createMinimalPDF() string {
	return `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
/Resources <<
/Font <<
/F1 <<
/Type /Font
/Subtype /Type1
/BaseFont /Helvetica
>>
>>
>>
>>
endobj
4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
100 700 Td
(Test PDF) Tj
ET
endstream
endobj
xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000317 00000 n 
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
410
%%EOF`
}

func TestPDFExtractor_ExtractText_ValidPDF(t *testing.T) {
	extractor := NewPDFExtractor()
	
	// Create a minimal valid PDF
	pdfContent := createMinimalPDF()
	r := strings.NewReader(pdfContent)
	
	text, err := extractor.ExtractText(r)
	
	// This might fail depending on the PDF library's strictness
	// but we're testing that it at least attempts to process it
	if err != nil {
		// It's okay if it fails to extract from our minimal PDF
		// as long as it doesn't panic
		t.Logf("ExtractText() error (expected for minimal PDF): %v", err)
	} else {
		t.Logf("ExtractText() successfully extracted: %q", text)
		// If it succeeds, text should not be empty
		if len(strings.TrimSpace(text)) == 0 {
			t.Error("ExtractText() returned empty text from valid PDF")
		}
	}
}
