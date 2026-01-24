package post

import (
	"errors"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
* PostRepo handles database operations for posts
 */
type PostRepo struct {
	db *gorm.DB
}

/*
* NewPostRepo creates a new PostRepo instance
*/
func NewPostRepo(db *gorm.DB) *PostRepo { return &PostRepo{db: db} }

/*
* Create inserts a new post into the database
*/
func (r *PostRepo) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

/*
* Update modifies an existing post in the database
*/
func (r *PostRepo) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

/*
* Delete removes a post from the database by ID
*/
func (r *PostRepo) Delete(id string) error {
	return r.db.Delete(&models.Post{}, "id = ?", id).Error
}

/*
* GetByID retrieves a post by its ID
*/
func (r *PostRepo) GetByID(id string) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, err
}

/*
* List retrieves all posts, ordered by creation time
*/
func (r *PostRepo) List() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Order("created_at desc").Find(&posts).Error
	return posts, err
}
