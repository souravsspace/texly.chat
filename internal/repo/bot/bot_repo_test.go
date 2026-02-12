package bot

import (
	"testing"

	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestBotRepo_Create(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	bot := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       "user-1",
		Name:         "Test Bot",
		SystemPrompt: "You are a test bot",
	}

	// Test Create
	err := repo.Create(bot)
	assert.NoError(t, err)

	// Verify existence
	var count int64
	testDB.Model(&models.Bot{}).Where("id = ?", bot.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	// Cleanup
	testDB.Unscoped().Delete(bot)
}

func TestBotRepo_GetByUserID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	userID := "user-1"
	bot1 := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         "Bot 1",
		SystemPrompt: "Prompt 1",
	}
	bot2 := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         "Bot 2",
		SystemPrompt: "Prompt 2",
	}
	otherBot := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       "other-user",
		Name:         "Other Bot",
		SystemPrompt: "Other Prompt",
	}

	repo.Create(bot1)
	repo.Create(bot2)
	repo.Create(otherBot)

	// Cleanup at end
	defer func() {
		testDB.Unscoped().Delete(bot1)
		testDB.Unscoped().Delete(bot2)
		testDB.Unscoped().Delete(otherBot)
	}()

	// Test GetByUserID
	bots, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Len(t, bots, 2)

	// Verify contents (names should match)
	names := []string{bots[0].Name, bots[1].Name}
	assert.Contains(t, names, "Bot 1")
	assert.Contains(t, names, "Bot 2")
}

func TestBotRepo_GetByID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	userID := "user-1"
	bot := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         "Test Bot",
		SystemPrompt: "Prompt",
	}
	repo.Create(bot)

	// Test Found
	found, err := repo.GetByID(bot.ID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, bot.Name, found.Name)

	// Test Not Found (Wrong ID)
	notFound, err := repo.GetByID("wrong-id", userID)
	assert.NoError(t, err)
	assert.Nil(t, notFound)

	// Test Not Found (Wrong UserID)
	notFoundUser, err := repo.GetByID(bot.ID, "wrong-user")
	assert.NoError(t, err)
	assert.Nil(t, notFoundUser)
}

func TestBotRepo_Update(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	bot := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       "user-1",
		Name:         "Old Name",
		SystemPrompt: "Old Prompt",
	}
	repo.Create(bot)

	// Update
	bot.Name = "New Name"
	err := repo.Update(bot)
	assert.NoError(t, err)

	// Verify
	updated, err := repo.GetByID(bot.ID, bot.UserID)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
}

func TestBotRepo_Delete(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	userID := "user-1"
	bot := &models.Bot{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   "To Delete",
	}
	repo.Create(bot)

	// Test Delete
	err := repo.Delete(bot.ID, userID)
	assert.NoError(t, err)

	// Verify deleted
	found, err := repo.GetByID(bot.ID, userID)
	assert.NoError(t, err)
	assert.Nil(t, found)

	// Test Delete NonExistent
	err = repo.Delete("non-existent", userID)
	assert.Error(t, err)
}

func TestBotRepo_GetByIDPublic(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewBotRepo(testDB, nil)

	userID := "user-1"
	bot := &models.Bot{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         "Public Test Bot",
		SystemPrompt: "Public Prompt",
	}
	repo.Create(bot)
	defer testDB.Unscoped().Delete(bot)

	// Test Found (without user authentication)
	found, err := repo.GetByIDPublic(bot.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, bot.Name, found.Name)
	assert.Equal(t, bot.UserID, found.UserID)

	// Test Not Found
	notFound, err := repo.GetByIDPublic("wrong-id")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}
