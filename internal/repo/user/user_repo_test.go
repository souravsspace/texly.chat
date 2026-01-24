package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewUserRepo(testDB)
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	/*
	* Test Create
	 */
	err := repo.Create(user)
	assert.NoError(t, err)

	/*
	* Verify existence
	 */
	var count int64
	testDB.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	/*
	* Cleanup
	 */
	testDB.Delete(user)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	// Setup
	testDB := shared.SetupTestDB()
	repo := NewUserRepo(testDB)
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "getbyemail@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Get By Email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.Create(user)
	defer testDB.Delete(user)

	/*
	* Test Found
	 */
	found, err := repo.GetByEmail("getbyemail@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)

	/*
	* Test Not Found
	 */
	notFound, err := repo.GetByEmail("nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestUserRepo_GetByID(t *testing.T) {
	// Setup
	testDB := shared.SetupTestDB()
	repo := NewUserRepo(testDB)
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "getbyid@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Get By ID",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.Create(user)
	defer testDB.Delete(user)

	// Test Found
	found, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Email, found.Email)

	// Test Not Found
	notFound, err := repo.GetByID("nonexistent-id")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}
