package extractor

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

/*
* ExcelParser handles Excel and CSV file parsing
 */
type ExcelParser struct{}

/*
* NewExcelParser creates a new Excel parser instance
 */
func NewExcelParser() *ExcelParser {
	return &ExcelParser{}
}

/*
* ParseExcel parses an Excel file and returns text representation
 */
func (p *ExcelParser) ParseExcel(reader io.Reader) (string, error) {
	// Open Excel file from reader
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	var result strings.Builder

	// Get all sheet names
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("no sheets found in Excel file")
	}

	// Process each sheet
	for _, sheetName := range sheets {
		result.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheetName))

		// Get all rows in the sheet
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return "", fmt.Errorf("failed to get rows from sheet %s: %w", sheetName, err)
		}

		// Convert rows to text
		for _, row := range rows {
			result.WriteString(strings.Join(row, "\t"))
			result.WriteString("\n")
		}

		result.WriteString("\n")
	}

	text := result.String()
	if len(text) == 0 {
		return "", fmt.Errorf("no data could be extracted from Excel file")
	}

	return text, nil
}

/*
* ParseCSV parses a CSV file and returns text representation
 */
func (p *ExcelParser) ParseCSV(reader io.Reader) (string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := csvReader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no data found in CSV file")
	}

	var result strings.Builder

	// Convert CSV records to text
	for _, record := range records {
		result.WriteString(strings.Join(record, "\t"))
		result.WriteString("\n")
	}

	return result.String(), nil
}
