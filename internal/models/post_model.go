package models

import "time"

/*
* Post represents a user-created post
 */
type Post struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
* CreatePostRequest holds data for creating a new post
*/
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

/*
* UpdatePostRequest holds data for updating an existing post
*/
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
