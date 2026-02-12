package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/services/oauth"
)

type GoogleHandler struct {
	oauthService *oauth.OAuthService
	stateService *oauth.StateService
	cfg          configs.Config
}

func NewGoogleHandler(oauthService *oauth.OAuthService, stateService *oauth.StateService, cfg configs.Config) *GoogleHandler {
	return &GoogleHandler{
		oauthService: oauthService,
		stateService: stateService,
		cfg:          cfg,
	}
}

func (h *GoogleHandler) GoogleLogin(c *gin.Context) {
	state, err := h.stateService.GenerateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	url := h.oauthService.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *GoogleHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state or code"})
		return
	}

	// Validate state
	if !h.stateService.ValidateState(state) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired state"})
		return
	}

	// Exchange code for user
	user, err := h.oauthService.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, "oauth_failed"))
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, "token_generation_failed"))
		return
	}

	// Redirect to frontend with token
	// Using fragment to avoid sending token to server in subsequent requests in URL
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback#token=%s", h.cfg.FrontendURL, tokenString))
}
