package storage

import (
	"testing"
)

func TestValidateFileType(t *testing.T) {
	svc := &MinIOStorageService{
		allowedTypes: map[string]bool{
			".txt":  true,
			".md":   true,
			".pdf":  true,
			".xlsx": true,
			".xls":  true,
			".csv":  true,
		},
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"valid txt", "test.txt", false},
		{"valid md", "README.md", false},
		{"valid pdf", "document.pdf", false},
		{"valid xlsx", "spreadsheet.xlsx", false},
		{"valid xls", "old-spreadsheet.xls", false},
		{"valid csv", "data.csv", false},
		{"invalid exe", "virus.exe", true},
		{"invalid zip", "archive.zip", true},
		{"invalid jpg", "image.jpg", true},
		{"uppercase extension", "file.TXT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateFileType(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	svc := &MinIOStorageService{
		maxUploadSizeMB: 100,
	}

	tests := []struct {
		name    string
		size    int64
		wantErr bool
	}{
		{"valid 1MB", 1 * 1024 * 1024, false},
		{"valid 50MB", 50 * 1024 * 1024, false},
		{"valid 100MB", 100 * 1024 * 1024, false},
		{"invalid 101MB", 101 * 1024 * 1024, true},
		{"invalid 200MB", 200 * 1024 * 1024, true},
		{"invalid 0 bytes", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateFileSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateObjectName(t *testing.T) {
	svc := &MinIOStorageService{}

	tests := []struct {
		name             string
		sourceID         string
		originalFilename string
		want             string
	}{
		{"simple filename", "123", "test.txt", "sources/123/test.txt"},
		{"filename with spaces", "456", "my document.pdf", "sources/456/my document.pdf"},
		{"filename with special chars", "789", "file-name_v2.xlsx", "sources/789/file-name_v2.xlsx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.GenerateObjectName(tt.sourceID, tt.originalFilename)
			if got != tt.want {
				t.Errorf("GenerateObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{"txt file", "test.txt", "text/plain"},
		{"md file", "README.md", "text/markdown"},
		{"pdf file", "doc.pdf", "application/pdf"},
		{"xlsx file", "sheet.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		{"xls file", "old.xls", "application/vnd.ms-excel"},
		{"csv file", "data.csv", "text/csv"},
		{"unknown file", "file.unknown", "application/octet-stream"},
		{"uppercase extension", "file.PDF", "application/pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContentType(tt.filename)
			if got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
