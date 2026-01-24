package post

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestPostRepo_Create(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewPostRepo(testDB)
	post := &models.Post{
		ID:        uuid.New().String(),
		Content:   "Test Post Content",
		UserID:    uuid.New().String(),
		/*
		* Usually foreign key, but sqlite might be chill if FK generic not enforced or user doesn't exist yet? Gorm enforces FK if defined.
		 */
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	/*
	* Create a dummy user first to satisfy FK constraint if needed
	 */
	user := &models.User{
		ID:           post.UserID,
		Email:        "postwriter@example.com",
		PasswordHash: "hash",
		Name:         "Writer",
	}
	testDB.Create(user)
	defer testDB.Delete(user)

	/*
	* Test Create
	 */
	err := repo.Create(post)
	assert.NoError(t, err)

	/*
	* Verify
	 */
	var count int64
	testDB.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	/*
	* Cleanup
	 */
	testDB.Delete(post)
}

func TestPostRepo_GetAll(t *testing.T) {
	// Setup
	testDB := shared.SetupTestDB()
	repo := NewPostRepo(testDB)
	
	// Create user
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "viewer@example.com", Name: "Viewer"}
	testDB.Create(user)
	defer testDB.Delete(user)

	/*
	* Create posts
	 */
	post1 := &models.Post{ID: uuid.New().String(), Content: "Post 1", UserID: userID}
	post2 := &models.Post{ID: uuid.New().String(), Content: "Post 2", UserID: userID}
	repo.Create(post1)
	repo.Create(post2)
	defer testDB.Delete(post1)
	defer testDB.Delete(post2)

	/*
	* Test List
	 */
	posts, err := repo.List()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(posts), 2)
}

func TestPostRepo_GetByID(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewPostRepo(testDB)
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "getter@example.com", Name: "Getter"}
	testDB.Create(user)
	defer testDB.Delete(user)

	post := &models.Post{ID: uuid.New().String(), Content: "Target Post", UserID: userID}
	repo.Create(post)
	defer testDB.Delete(post)

	/*
	* Found
	 */
	found, err := repo.GetByID(post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, post.Content, found.Content)

	/*
	* Not Found
	 */
	notFound, err := repo.GetByID("missing")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestPostRepo_Update(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewPostRepo(testDB)
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "updater@example.com", Name: "Updater"}
	testDB.Create(user)
	defer testDB.Delete(user)

	post := &models.Post{ID: uuid.New().String(), Content: "Old Content", UserID: userID}
	repo.Create(post)
	defer testDB.Delete(post)

	/*
	* Test Update
	 */
	post.Content = "New Content"
	err := repo.Update(post)
	assert.NoError(t, err)

	/*
	* Verify
	 */
	updated, _ := repo.GetByID(post.ID)
	assert.Equal(t, "New Content", updated.Content)
}

func TestPostRepo_Delete(t *testing.T) {
	testDB := shared.SetupTestDB()
	repo := NewPostRepo(testDB)
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "deleter@example.com", Name: "Deleter"}
	testDB.Create(user)
	defer testDB.Delete(user)

	post := &models.Post{ID: uuid.New().String(), Content: "Delete Me", UserID: userID}
	repo.Create(post)

	/*
	* Test Delete
	 */
	err := repo.Delete(post.ID)
	assert.NoError(t, err)

	/*
	* Verify
	 */
	found, _ := repo.GetByID(post.ID)
	/*
	* As per repo implementation, returns nil if not found
	 */
	assert.Nil(t, found)
}
