package models

import "time"

/*
* User represents a registered user in the system
 */
type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-"` // Nullable for OAuth users, removed gorm:not null
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	GoogleID     *string   `json:"google_id" gorm:"unique"`
	AuthProvider string    `json:"auth_provider" gorm:"default:'email'"` // email, google
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

/*
* LoginRequest holds the credentials for user login
 */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
* SignupRequest holds data for creating a new user
 */
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

/*
* AuthResponse is the response payload for successful authentication
 */
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
