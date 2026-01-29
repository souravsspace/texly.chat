package extractor

import (
	"strings"
	"testing"
)

func TestExcelParser_ParseCSV(t *testing.T) {
	parser := NewExcelParser()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid CSV with headers",
			content: "Name,Age,City\nJohn,30,New York\nJane,25,Los Angeles",
			wantErr: false,
		},
		{
			name:    "CSV with quotes",
			content: "Name,Description\n\"John Doe\",\"Software Engineer\"\n\"Jane Smith\",\"Data Scientist\"",
			wantErr: false,
		},
		{
			name:    "CSV with empty fields",
			content: "Name,Age,City\nJohn,,New York\n,25,",
			wantErr: false,
		},
		{
			name:    "single row CSV",
			content: "Header1,Header2,Header3",
			wantErr: false,
		},
		{
			name:    "empty CSV",
			content: "",
			wantErr: true,
		},
		{
			name:    "CSV with variable columns",
			content: "A,B,C\n1,2\n3,4,5,6",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.content)
			text, err := parser.ParseCSV(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(text) == 0 {
					t.Error("ParseCSV() returned empty text for valid content")
				}
				// Check that text contains tab-separated values
				if !strings.Contains(text, "\t") && len(strings.Split(tt.content, ",")) > 1 {
					t.Error("ParseCSV() did not convert commas to tabs")
				}
			}
		})
	}
}

func TestExcelParser_ParseCSV_MultipleRows(t *testing.T) {
	parser := NewExcelParser()
	
	csv := `Product,Price,Quantity
Apple,1.20,50
Banana,0.50,100
Orange,0.80,75`
	
	r := strings.NewReader(csv)
	text, err := parser.ParseCSV(r)
	
	if err != nil {
		t.Errorf("ParseCSV() failed: %v", err)
	}
	
	// Should have 4 lines (3 rows + newlines)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) != 4 {
		t.Errorf("ParseCSV() expected 4 lines, got %d", len(lines))
	}
	
	// First line should contain tab-separated headers
	if !strings.Contains(lines[0], "Product\tPrice\tQuantity") {
		t.Error("ParseCSV() did not properly format headers")
	}
}

// Note: Testing ParseExcel requires creating actual Excel files or using mocks
// For now, we'll add a basic structure test
func TestExcelParser_ParseExcel_InvalidInput(t *testing.T) {
	parser := NewExcelParser()
	
	// Invalid Excel content (just plain text)
	r := strings.NewReader("This is not an Excel file")
	_, err := parser.ParseExcel(r)
	
	if err == nil {
		t.Error("ParseExcel() should fail on invalid Excel content")
	}
}
