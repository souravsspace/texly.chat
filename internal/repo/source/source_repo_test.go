package source

import (
	"testing"

	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestSourceRepo_Create(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	source := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusPending,
	}

	// Test Create
	err := repo.Create(source)
	assert.NoError(t, err)

	// Verify existence
	var count int64
	testDB.Model(&models.Source{}).Where("id = ?", source.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	// Cleanup
	testDB.Unscoped().Delete(source)
}

func TestSourceRepo_GetByID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	source := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusPending,
	}
	repo.Create(source)

	defer testDB.Unscoped().Delete(source)

	// Test Found
	found, err := repo.GetByID(source.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, source.URL, found.URL)

	// Test Not Found
	notFound, err := repo.GetByID("wrong-id")
	assert.Error(t, err)
	assert.Nil(t, notFound)
}

func TestSourceRepo_ListByBotID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	botID := "bot-1"
	source1 := &models.Source{
		ID:     uuid.New().String(),
		BotID:  botID,
		URL:    "https://example1.com",
		Status: models.SourceStatusPending,
	}
	source2 := &models.Source{
		ID:     uuid.New().String(),
		BotID:  botID,
		URL:    "https://example2.com",
		Status: models.SourceStatusCompleted,
	}
	otherSource := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "other-bot",
		URL:    "https://other.com",
		Status: models.SourceStatusPending,
	}

	repo.Create(source1)
	repo.Create(source2)
	repo.Create(otherSource)

	defer func() {
		testDB.Unscoped().Delete(source1)
		testDB.Unscoped().Delete(source2)
		testDB.Unscoped().Delete(otherSource)
	}()

	// Test ListByBotID
	sources, err := repo.ListByBotID(botID)
	assert.NoError(t, err)
	assert.Len(t, sources, 2)

	// Verify contents
	urls := []string{sources[0].URL, sources[1].URL}
	assert.Contains(t, urls, "https://example1.com")
	assert.Contains(t, urls, "https://example2.com")
}

func TestSourceRepo_UpdateStatus(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	source := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusPending,
	}
	repo.Create(source)

	defer testDB.Unscoped().Delete(source)

	// Update to processing
	err := repo.UpdateStatus(source.ID, models.SourceStatusProcessing, "")
	assert.NoError(t, err)

	// Verify
	updated, err := repo.GetByID(source.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.SourceStatusProcessing, updated.Status)
	assert.Nil(t, updated.ProcessedAt)

	// Update to completed (should set processed_at)
	err = repo.UpdateStatus(source.ID, models.SourceStatusCompleted, "")
	assert.NoError(t, err)

	completed, err := repo.GetByID(source.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.SourceStatusCompleted, completed.Status)
	assert.NotNil(t, completed.ProcessedAt)

	// Update to failed with error message
	errorMsg := "Failed to fetch"
	err = repo.UpdateStatus(source.ID, models.SourceStatusFailed, errorMsg)
	assert.NoError(t, err)

	failed, err := repo.GetByID(source.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.SourceStatusFailed, failed.Status)
	assert.Equal(t, errorMsg, failed.ErrorMessage)
	assert.NotNil(t, failed.ProcessedAt)
}

func TestSourceRepo_Delete(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	source := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusPending,
	}
	repo.Create(source)

	// Test Delete
	err := repo.Delete(source.ID)
	assert.NoError(t, err)

	// Verify soft deleted (still in DB but with DeletedAt set)
	var count int64
	testDB.Unscoped().Model(&models.Source{}).Where("id = ?", source.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	// Verify not found in normal queries
	_, err = repo.GetByID(source.ID)
	assert.Error(t, err)

	// Cleanup
	testDB.Unscoped().Delete(source)
}

func TestSourceRepo_GetByBotIDAndSourceID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewSourceRepo(testDB, nil)

	source := &models.Source{
		ID:     uuid.New().String(),
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusPending,
	}
	repo.Create(source)

	defer testDB.Unscoped().Delete(source)

	// Test Found with correct bot ID
	found, err := repo.GetByBotIDAndSourceID("bot-1", source.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, source.URL, found.URL)

	// Test Not Found with wrong bot ID
	notFound, err := repo.GetByBotIDAndSourceID("wrong-bot", source.ID)
	assert.Error(t, err)
	assert.Nil(t, notFound)
	assert.Contains(t, err.Error(), "not found or access denied")
}
