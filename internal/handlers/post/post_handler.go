package post

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	repo "github.com/souravsspace/texly.chat/internal/repo/post"
)

/*
* PostHandler handles post-related requets
 */
type PostHandler struct {
	repo *repo.PostRepo
}

/*
* NewPostHandler creates a new PostHandler instance
*/
func NewPostHandler(pr *repo.PostRepo) *PostHandler { return &PostHandler{repo: pr} }

/*
* CreatePost handles creation of a new post
*/
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	post := &models.Post{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.Create(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

/*
* GetPosts handles retrieval of all posts
*/
func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

/*
* GetPost handles retrieval of a single post by ID
*/
func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	post, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

/*
* UpdatePost handles updating an existing post
*/
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	post, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	if post.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.UpdatedAt = time.Now()

	if err := h.repo.Update(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

/*
* DeletePost handles deletion of a post
*/
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	post, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	if post.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
