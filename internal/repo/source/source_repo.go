package source

import (
	"context"
	"fmt"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/cache"
	"gorm.io/gorm"
)

/*
* SourceRepo handles database operations for sources
 */
type SourceRepo struct {
	db    *gorm.DB
	cache *cache.CacheService
}

/*
* NewSourceRepo creates a new source repository
 */
func NewSourceRepo(db *gorm.DB, cache *cache.CacheService) *SourceRepo {
	return &SourceRepo{db: db, cache: cache}
}

/*
* Create creates a new source
 */
func (r *SourceRepo) Create(source *models.Source) error {
	if err := r.db.Create(source).Error; err != nil {
		return err
	}
	// Invalidate source list cache for this bot
	_ = r.cache.DeletePattern(context.Background(), fmt.Sprintf(cache.SourceListCacheKey, source.BotID))
	return nil
}

/*
* GetByID retrieves a source by ID
 */
func (r *SourceRepo) GetByID(id string) (*models.Source, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.SourceCacheKey, id)

	// Try cache first
	var source models.Source
	if err := r.cache.GetJSON(ctx, cacheKey, &source); err == nil {
		return &source, nil
	}

	// Cache miss - query database
	if err := r.db.First(&source, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, source, cache.SourceCacheTTL)
	return &source, nil
}

/*
* ListByBotID retrieves all sources for a bot
 */
func (r *SourceRepo) ListByBotID(botID string) ([]*models.Source, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.SourceListCacheKey, botID)

	// Try cache first
	var sources []*models.Source
	if err := r.cache.GetJSON(ctx, cacheKey, &sources); err == nil {
		return sources, nil
	}

	// Cache miss - query database
	if err := r.db.Where("bot_id = ?", botID).Order("created_at DESC").Find(&sources).Error; err != nil {
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, sources, cache.SourceListCacheTTL)
	return sources, nil
}

/*
* UpdateStatus updates the status of a source
 */
func (r *SourceRepo) UpdateStatus(id string, status models.SourceStatus, errorMsg string) error {
	updates := map[string]interface{}{
		"status":        status,
		"error_message": errorMsg,
	}

	if status == models.SourceStatusCompleted || status == models.SourceStatusFailed {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if err := r.db.Model(&models.Source{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	// Invalidate source cache
	_ = r.cache.Delete(context.Background(), fmt.Sprintf(cache.SourceCacheKey, id))
	return nil
}

/*
* Delete soft deletes a source
 */
func (r *SourceRepo) Delete(id string) error {
	if err := r.db.Delete(&models.Source{}, "id = ?", id).Error; err != nil {
		return err
	}
	// Invalidate source cache
	_ = r.cache.Delete(context.Background(), fmt.Sprintf(cache.SourceCacheKey, id))
	return nil
}

/*
* GetByBotIDAndSourceID retrieves a source by bot ID and source ID (for authorization)
 */
func (r *SourceRepo) GetByBotIDAndSourceID(botID, sourceID string) (*models.Source, error) {
	var source models.Source
	if err := r.db.Where("bot_id = ? AND id = ?", botID, sourceID).First(&source).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("source not found or access denied")
		}
		return nil, err
	}
	return &source, nil
}

/*
* UpdateProgress updates the processing progress of a source (0-100)
 */
func (r *SourceRepo) UpdateProgress(id string, progress int) error {
	return r.db.Model(&models.Source{}).Where("id = ?", id).Update("processing_progress", progress).Error
}

/*
* UpdateFilePath updates the file path of a source
 */
func (r *SourceRepo) UpdateFilePath(id string, filePath string) error {
	return r.db.Model(&models.Source{}).Where("id = ?", id).Update("file_path", filePath).Error
}
