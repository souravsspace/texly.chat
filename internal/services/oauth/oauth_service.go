package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
)

type OAuthService struct {
	googleConfig *oauth2.Config
	db           *gorm.DB
	frontendURL  string
}

func NewOAuthService(cfg configs.Config, db *gorm.DB) *OAuthService {
	return &OAuthService{
		googleConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		db:          db,
		frontendURL: cfg.FrontendURL,
	}
}

func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.googleConfig.AuthCodeURL(state)
}

func (s *OAuthService) HandleGoogleCallback(ctx context.Context, code string) (*models.User, error) {
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.googleConfig.Client(ctx, token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer userInfo.Body.Close()

	// Parse user info
	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(userInfo.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	if !googleUser.VerifiedEmail {
		return nil, errors.New("email not verified by Google")
	}

	// Find or create user
	var user models.User
	result := s.db.Where("email = ?", googleUser.Email).First(&user)

	if result.Error == nil {
		// User exists, update Google ID if not set
		if user.GoogleID == nil {
			user.GoogleID = &googleUser.ID
			user.AuthProvider = "google"
			if user.Avatar == "" {
				user.Avatar = googleUser.Picture
			}
			if err := s.db.Save(&user).Error; err != nil {
				return nil, err
			}
		}
		return &user, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new user
		newUser := models.User{
			ID:           uuid.NewString(), // assuming uuid is imported or available via google/uuid
			Email:        googleUser.Email,
			Name:         googleUser.Name,
			Avatar:       googleUser.Picture,
			GoogleID:     &googleUser.ID,
			AuthProvider: "google",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.db.Create(&newUser).Error; err != nil {
			return nil, err
		}
		return &newUser, nil
	} else {
		return nil, result.Error
	}
}
