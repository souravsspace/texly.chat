package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/*
* MinIOStorageService handles file storage operations with MinIO
 */
type MinIOStorageService struct {
	client          *minio.Client
	bucket          string
	maxUploadSizeMB int
	allowedTypes    map[string]bool
}

/*
* NewMinIOStorageService creates a new MinIO storage service instance
 */
func NewMinIOStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool, maxUploadSizeMB int) (*MinIOStorageService, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists, if not we'll let the docker init handle it
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket '%s' does not exist - please run MinIO initialization", bucket)
	}

	// Define allowed file types
	allowedTypes := map[string]bool{
		".txt":  true,
		".md":   true,
		".pdf":  true,
		".xlsx": true,
		".xls":  true,
		".csv":  true,
	}

	return &MinIOStorageService{
		client:          minioClient,
		bucket:          bucket,
		maxUploadSizeMB: maxUploadSizeMB,
		allowedTypes:    allowedTypes,
	}, nil
}

/*
* UploadFile uploads a file to MinIO
* objectName: path in MinIO (e.g., "sources/source-id/filename.pdf")
* reader: file content reader
* size: file size in bytes
* contentType: MIME type
 */
func (s *MinIOStorageService) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	// Upload to MinIO
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return nil
}

/*
* GetFile retrieves a file from MinIO
 */
func (s *MinIOStorageService) GetFile(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from MinIO: %w", err)
	}
	return object, nil
}

/*
* DeleteFile deletes a file from MinIO
 */
func (s *MinIOStorageService) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}
	return nil
}

/*
* ValidateFileType checks if the file extension is allowed
 */
func (s *MinIOStorageService) ValidateFileType(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !s.allowedTypes[ext] {
		return fmt.Errorf("file type '%s' is not supported. Allowed types: .txt, .md, .pdf, .xlsx, .xls, .csv", ext)
	}
	return nil
}

/*
* ValidateFileSize checks if the file size is within limits
 */
func (s *MinIOStorageService) ValidateFileSize(size int64) error {
	maxBytes := int64(s.maxUploadSizeMB) * 1024 * 1024
	if size > maxBytes {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d MB)", size, s.maxUploadSizeMB)
	}
	if size == 0 {
		return fmt.Errorf("file size cannot be 0 bytes")
	}
	return nil
}

/*
* GenerateObjectName creates a unique object path in MinIO
* Format: sources/{sourceID}/{originalFilename}
 */
func (s *MinIOStorageService) GenerateObjectName(sourceID, originalFilename string) string {
	return fmt.Sprintf("sources/%s/%s", sourceID, originalFilename)
}

/*
* GetContentType returns the MIME type based on file extension
 */
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	contentTypes := map[string]string{
		".txt":  "text/plain",
		".md":   "text/markdown",
		".pdf":  "application/pdf",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".xls":  "application/vnd.ms-excel",
		".csv":  "text/csv",
	}
	
	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
