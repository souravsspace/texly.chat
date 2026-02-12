package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/souravsspace/texly.chat/internal/models"
	repo "github.com/souravsspace/texly.chat/internal/repo/user"
	"github.com/souravsspace/texly.chat/internal/services/oauth"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func SetupTestRedis() (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestAuthHandler_Signup(t *testing.T) {
	testDB := shared.SetupTestDB()
	testCfg := shared.GetTestConfig()
	userRepo := repo.NewUserRepo(testDB)
	handler := NewAuthHandler(userRepo, testCfg)
	router := gin.New()
	router.POST("/auth/signup", handler.Signup)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.SignupRequest{
			Email:    "newuser@example.com",
			Password: "password123",
			Name:     "New User",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, reqBody.Email, response.User.Email)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		existing := models.SignupRequest{
			Email:    "duplicate@example.com",
			Password: "password123",
			Name:     "Existing User",
		}
		jsonBody, _ := json.Marshal(existing)
		req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	testDB := shared.SetupTestDB()
	testCfg := shared.GetTestConfig()
	userRepo := repo.NewUserRepo(testDB)
	handler := NewAuthHandler(userRepo, testCfg)
	router := gin.New()
	router.POST("/auth/signup", handler.Signup)
	router.POST("/auth/login", handler.Login)

	userReq := models.SignupRequest{
		Email:    "loginuser@example.com",
		Password: "password123",
		Name:     "Login User",
	}
	jsonBody, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(httptest.NewRecorder(), req)

	t.Run("Success", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "loginuser@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "loginuser@example.com",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGoogleHandler_GoogleLogin(t *testing.T) {
	testDB := shared.SetupTestDB()
	testCfg := shared.GetTestConfig()

	client, mr := SetupTestRedis()
	defer mr.Close()

	stateService := oauth.NewStateService(client)
	oauthService := oauth.NewOAuthService(testCfg, testDB)

	handler := NewGoogleHandler(oauthService, stateService, testCfg)

	router := gin.New()
	router.GET("/auth/google", handler.GoogleLogin)

	t.Run("Redirects to Google", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/auth/google", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		assert.Contains(t, location, "accounts.google.com")
		assert.Contains(t, location, "state=")
	})
}
