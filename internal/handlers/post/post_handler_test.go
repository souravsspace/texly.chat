package post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	repo "github.com/souravsspace/texly.chat/internal/repo/post"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler_CreatePost(t *testing.T) {
	testDB := shared.SetupTestDB()
	postRepo := repo.NewPostRepo(testDB)
	handler := NewPostHandler(postRepo)
	
	/*
	* Create user for FK
	 */
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "poster@example.com", Name: "Poster"}
	testDB.Create(user)
	defer testDB.Delete(user)

	router := gin.New()
	/*
	* Mock Auth Middleware
	 */
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/posts", handler.CreatePost)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreatePostRequest{
			Title:   "New Post",
			Content: "This is a new post",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var post models.Post
		json.Unmarshal(w.Body.Bytes(), &post)
		assert.Equal(t, reqBody.Title, post.Title)
		assert.Equal(t, userID, post.UserID)
	})
}

func TestPostHandler_GetPosts(t *testing.T) {
	testDB := shared.SetupTestDB()
	postRepo := repo.NewPostRepo(testDB)
	handler := NewPostHandler(postRepo)
	router := gin.New()
	router.GET("/posts", handler.GetPosts)

	/*
	* Seed data
	 */
	userID := uuid.New().String()
	user := &models.User{ID: userID, Email: "list@example.com", Name: "Lister"}
	testDB.Create(user)
	testDB.Create(&models.Post{ID: uuid.New().String(), UserID: userID, Title: "P1", Content: "C1"})
	
	req, _ := http.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var posts []models.Post
	json.Unmarshal(w.Body.Bytes(), &posts)
	assert.NotEmpty(t, posts)
}

func TestPostHandler_UpdatePost(t *testing.T) {
	testDB := shared.SetupTestDB()
	postRepo := repo.NewPostRepo(testDB)
	handler := NewPostHandler(postRepo)

	userID := uuid.New().String()
	otherUserID := uuid.New().String()
	
	/*
	* Setup users
	 */
	testDB.Create(&models.User{ID: userID, Email: "updater1@example.com", Name: "U1"})
	testDB.Create(&models.User{ID: otherUserID, Email: "updater2@example.com", Name: "U2"})

	/*
	* Setup post
	 */
	post := &models.Post{ID: uuid.New().String(), UserID: userID, Title: "Original", Content: "Original"}
	testDB.Create(post)

	router := gin.New()
	router.PUT("/posts/:id", func(c *gin.Context) {
		/*
		* Middleware to extract user_id from header for testing different users
		 */
		uid := c.GetHeader("X-User-ID")
		c.Set("user_id", uid)
		c.Next()
	}, handler.UpdatePost)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.UpdatePostRequest{Title: "Updated", Content: "Updated"}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/posts/"+post.ID, bytes.NewBuffer(jsonBody))
		req.Header.Set("X-User-ID", userID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		
		var updated models.Post
		testDB.First(&updated, "id = ?", post.ID)
		assert.Equal(t, "Updated", updated.Title)
	})

	t.Run("Forbidden", func(t *testing.T) {
		reqBody := models.UpdatePostRequest{Title: "Hacked", Content: "Hacked"}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/posts/"+post.ID, bytes.NewBuffer(jsonBody))
		req.Header.Set("X-User-ID", otherUserID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestPostHandler_DeletePost(t *testing.T) {
	testDB := shared.SetupTestDB()
	postRepo := repo.NewPostRepo(testDB)
	handler := NewPostHandler(postRepo)

	userID := uuid.New().String()
	testDB.Create(&models.User{ID: userID, Email: "deleter1@example.com", Name: "D1"})
	
	post := &models.Post{ID: uuid.New().String(), UserID: userID, Title: "To Delete", Content: "Bye", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	testDB.Create(post)

	router := gin.New()
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, handler.DeletePost)

	req, _ := http.NewRequest("DELETE", "/posts/"+post.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	
	var count int64
	testDB.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
