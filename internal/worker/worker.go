package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/queue"
	sourceRepo "github.com/souravsspace/texly.chat/internal/repo/source"
	vectorRepo "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/chunker"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/scraper"
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
	maxChunkTokens int
}

/*
* NewWorker creates a new worker instance
 */
func NewWorker(
	db *gorm.DB,
	embeddingSvc *embedding.EmbeddingService,
	vectorRepo *vectorRepo.VectorRepository,
) *Worker {
	return &Worker{
		db:             db,
		sourceRepo:     sourceRepo.NewSourceRepo(db),
		vectorRepo:     vectorRepo,
		scraperSvc:     scraper.NewScraperService(),
		embeddingSvc:   embeddingSvc,
		maxChunkTokens: 800, // ~600 words per chunk
	}
}

/*
* ProcessScrapeJob is the handler function for scraping jobs
 */
func (w *Worker) ProcessScrapeJob(job queue.Job) error {
	fmt.Printf("Processing scrape job for source %s (URL: %s)\n", job.SourceID, job.URL)

	// Update status to processing
	if err := w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to update source status to processing: %w", err)
	}

	// Scrape the URL
	content, err := w.scraperSvc.FetchAndClean(job.URL)
	if err != nil {
		// Update status to failed
		errMsg := fmt.Sprintf("Failed to scrape: %v", err)
		_ = w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusFailed, errMsg)
		return fmt.Errorf("scraping failed: %w", err)
	}

	// Chunk the content
	chunks := chunker.ChunkText(content, w.maxChunkTokens)
	fmt.Printf("Created %d chunks from content\n", len(chunks))

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

	// Update status to completed
	if err := w.sourceRepo.UpdateStatus(job.SourceID, models.SourceStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update source status to completed: %w", err)
	}

	fmt.Printf("Successfully processed source %s: %d chunks created\n", job.SourceID, len(chunks))
	return nil
}
