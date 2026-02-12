package worker

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/queue"
	sourceRepo "github.com/souravsspace/texly.chat/internal/repo/source"
	vectorRepo "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/chunker"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/extractor"
	"github.com/souravsspace/texly.chat/internal/services/scraper"
	"github.com/souravsspace/texly.chat/internal/services/storage"
	"gorm.io/gorm"
)

/*
* Worker handles background job processing
 */
type Worker struct {
	db             *gorm.DB
	sourceRepo     *sourceRepo.SourceRepo
	vectorRepo     *vectorRepo.VectorRepository
	scraperSvc     *scraper.ScraperService
	embeddingSvc   *embedding.EmbeddingService
	storageSvc     *storage.MinIOStorageService
	pdfExtractor   *extractor.PDFExtractor
	excelParser    *extractor.ExcelParser
	textReader     *extractor.TextReader
	maxChunkTokens int
}

/*
* NewWorker creates a new worker instance
 */
func NewWorker(
	db *gorm.DB,
	embeddingSvc *embedding.EmbeddingService,
	vectorRepo *vectorRepo.VectorRepository,
	storageSvc *storage.MinIOStorageService,
	sourceRepoInstance *sourceRepo.SourceRepo,
) *Worker {
	return &Worker{
		db:             db,
		sourceRepo:     sourceRepoInstance,
		vectorRepo:     vectorRepo,
		scraperSvc:     scraper.NewScraperService(),
		embeddingSvc:   embeddingSvc,
		storageSvc:     storageSvc,
		pdfExtractor:   extractor.NewPDFExtractor(),
		excelParser:    extractor.NewExcelParser(),
		textReader:     extractor.NewTextReader(),
		maxChunkTokens: 800, // ~600 words per chunk
	}
}

/*
* ProcessScrapeJob is the handler function for processing jobs (scraping, file extraction, etc.)
 */
func (w *Worker) ProcessScrapeJob(job queue.Job) error {
	// Get source to determine type
	source, err := w.sourceRepo.GetByID(job.SourceID)
	if err != nil {
		return fmt.Errorf("failed to get source: %w", err)
	}

	fmt.Printf("Processing job for source %s (type: %s)\n", job.SourceID, source.SourceType)

	// Update status to processing
	if err := w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to update source status to processing: %w", err)
	}
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 10)

	// Extract content based on source type
	var content string
	switch source.SourceType {
	case models.SourceTypeURL:
		content, err = w.processURLSource(source)
	case models.SourceTypeFile:
		content, err = w.processFileSource(source)
	case models.SourceTypeText:
		content, err = w.processTextSource(source)
	default:
		err = fmt.Errorf("unknown source type: %s", source.SourceType)
	}

	if err != nil {
		errMsg := fmt.Sprintf("Failed to extract content: %v", err)
		_ = w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusFailed, errMsg)
		return err
	}
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 30)

	// Chunk the content
	chunks := chunker.ChunkText(content, w.maxChunkTokens)
	fmt.Printf("Created %d chunks from content\n", len(chunks))
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 50)

	// Save chunks to database first
	var savedChunks []models.DocumentChunk
	for i, chunkContent := range chunks {
		chunk := &models.DocumentChunk{
			SourceID:   job.SourceID,
			Content:    chunkContent,
			ChunkIndex: i,
			CreatedAt:  time.Now(),
		}

		if err := w.db.Create(chunk).Error; err != nil {
			errMsg := fmt.Sprintf("Failed to save chunk %d: %v", i, err)
			_ = w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusFailed, errMsg)
			return fmt.Errorf("failed to save chunk: %w", err)
		}

		savedChunks = append(savedChunks, *chunk)
	}
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 70)

	// Generate embeddings if embedding service is available
	if w.embeddingSvc != nil && w.vectorRepo != nil {
		ctx := context.Background()
		fmt.Printf("Generating embeddings for %d chunks...\n", len(savedChunks))

		embeddings, err := w.embeddingSvc.EmbedChunks(ctx, savedChunks)
		if err != nil {
			// Log error but don't fail the whole job
			errMsg := fmt.Sprintf("Warning: Failed to generate embeddings: %v", err)
			fmt.Println(errMsg)
			// Continue without embeddings - chunks are still searchable via full-text
		} else {
			// Store embeddings in vector database
			vectorData := make([]vectorRepo.VectorData, len(savedChunks))
			for i, chunk := range savedChunks {
				vectorData[i] = vectorRepo.VectorData{
					ChunkID:   chunk.ID,
					Embedding: embeddings[i],
				}
			}

			if err := w.vectorRepo.BulkInsertEmbeddings(ctx, vectorData); err != nil {
				// Log error but don't fail
				fmt.Printf("Warning: Failed to store embeddings: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully stored %d embeddings\n", len(vectorData))
			}
		}
	} else {
		fmt.Println("⚠️  Embedding service not configured - skipping vector embeddings")
	}
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 90)

	// Update status to completed
	if err := w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update source status to completed: %w", err)
	}
	_ = w.sourceRepo.UpdateProgress(job.SourceID, 100)

	fmt.Printf("Successfully processed source %s: %d chunks created\n", job.SourceID, len(chunks))
	return nil
}

/*
* processURLSource extracts content from a URL
 */
func (w *Worker) processURLSource(source *models.Source) (string, error) {
	fmt.Printf("Scraping URL: %s\n", source.URL)
	content, err := w.scraperSvc.FetchAndClean(source.URL)
	if err != nil {
		return "", fmt.Errorf("failed to scrape URL: %w", err)
	}
	return content, nil
}

/*
* processFileSource extracts content from a file stored in MinIO
 */
func (w *Worker) processFileSource(source *models.Source) (string, error) {
	fmt.Printf("Processing file: %s (type: %s)\n", source.OriginalFilename, source.ContentType)

	// Download file from MinIO
	ctx := context.Background()
	object, err := w.storageSvc.GetFile(ctx, source.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to download file from storage: %w", err)
	}
	defer object.Close()

	// Extract content based on file type
	ext := strings.ToLower(filepath.Ext(source.OriginalFilename))
	var content string

	switch ext {
	case ".pdf":
		content, err = w.pdfExtractor.ExtractText(object)
		if err != nil {
			return "", fmt.Errorf("failed to extract PDF content: %w", err)
		}
	case ".xlsx", ".xls":
		content, err = w.excelParser.ParseExcel(object)
		if err != nil {
			return "", fmt.Errorf("failed to parse Excel file: %w", err)
		}
	case ".csv":
		content, err = w.excelParser.ParseCSV(object)
		if err != nil {
			return "", fmt.Errorf("failed to parse CSV file: %w", err)
		}
	case ".txt", ".md":
		content, err = w.textReader.ReadTextFile(object)
		if err != nil {
			return "", fmt.Errorf("failed to read text file: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	return content, nil
}

/*
* processTextSource extracts content from a text source stored in MinIO
 */
func (w *Worker) processTextSource(source *models.Source) (string, error) {
	fmt.Printf("Processing text source: %s\n", source.OriginalFilename)

	// Download text file from MinIO
	ctx := context.Background()
	object, err := w.storageSvc.GetFile(ctx, source.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to download text from storage: %w", err)
	}
	defer object.Close()

	// Read text content
	content, err := w.textReader.ReadTextFile(object)
	if err != nil {
		return "", fmt.Errorf("failed to read text content: %w", err)
	}

	return content, nil
}
