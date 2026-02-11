package source

import (
	"fmt"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
* SourceRepo handles database operations for sources
 */
type SourceRepo struct {
	db *gorm.DB
}

/*
* NewSourceRepo creates a new source repository
 */
func NewSourceRepo(db *gorm.DB) *SourceRepo {
	return &SourceRepo{db: db}
}

/*
* Create creates a new source
 */
func (r *SourceRepo) Create(source *models.Source) error {
	return r.db.Create(source).Error
}

/*
* GetByID retrieves a source by ID
 */
func (r *SourceRepo) GetByID(id string) (*models.Source, error) {
	var source models.Source
	if err := r.db.First(&source, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

/*
* ListByBotID retrieves all sources for a bot
 */
func (r *SourceRepo) ListByBotID(botID string) ([]*models.Source, error) {
	var sources []*models.Source
	if err := r.db.Where("bot_id = ?", botID).Order("created_at DESC").Find(&sources).Error; err != nil {
		return nil, err
	}
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

	return r.db.Model(&models.Source{}).Where("id = ?", id).Updates(updates).Error
}

/*
* Delete soft deletes a source
 */
func (r *SourceRepo) Delete(id string) error {
	return r.db.Delete(&models.Source{}, "id = ?", id).Error
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
