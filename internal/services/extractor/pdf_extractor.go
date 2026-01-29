package extractor

import (
	"fmt"
	"io"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

/*
* PDFExtractor handles PDF text extraction
 */
type PDFExtractor struct{}

/*
* NewPDFExtractor creates a new PDF extractor instance
 */
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

/*
* ExtractText extracts plain text from a PDF reader
 */
func (e *PDFExtractor) ExtractText(reader io.ReadSeeker) (string, error) {
	// Create PDF reader
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Get number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get number ofpages: %w", err)
	}

	var allText string

	// Extract text from each page
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", fmt.Errorf("failed to get page %d: %w", i, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("failed to create extractor for page %d: %w", i, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}

		allText += text + "\n\n"
	}

	if len(allText) == 0 {
		return "", fmt.Errorf("no text could be extracted from PDF")
	}

	return allText, nil
}
