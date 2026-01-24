package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	repo "github.com/souravsspace/texly.chat/internal/repo/user"
)

/*
* UserHandler handles user-related requests
 */
type UserHandler struct {
	repo *repo.UserRepo
}

/*
* NewUserHandler creates a new UserHandler instance
*/
func NewUserHandler(ur *repo.UserRepo) *UserHandler { return &UserHandler{repo: ur} }

/*
* GetMe handles retrieval of the current authenticated user
*/
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	user, err := h.repo.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
