package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	repo "github.com/souravsspace/texly.chat/internal/repo/user"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetMe(t *testing.T) {
	testDB := shared.SetupTestDB()
	userRepo := repo.NewUserRepo(testDB)
	handler := NewUserHandler(userRepo)

	userID := uuid.New().String()
	user := &models.User{
		ID:           userID,
		Email:        "me@example.com",
		Name:         "Me",
		PasswordHash: "secret",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	testDB.Create(user)

	router := gin.New()
	router.GET("/users/me", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, handler.GetMe)

	req, _ := http.NewRequest("GET", "/users/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.User
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, user.Email, resp.Email)
	/*
	* Check json tag likely omits it, or handler might strip it? Default Gorm/JSON model usually includes unless ignored.
	 */
	assert.Empty(t, resp.PasswordHash)
	/*
	* Actually Handler.GetMe returns user from DB. Ideally PasswordHash should be json:"-" in mode.
	 */
}
