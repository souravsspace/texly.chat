package user

import (
	"errors"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
* UserRepo handles database operations for users
 */
type UserRepo struct {
	db *gorm.DB
}

/*
* NewUserRepo creates a new UserRepo instance
*/
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

/*
* Create inserts a new user into the database
*/
func (r *UserRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

/*
* GetByEmail retrieves a user by their email address
*/
func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

/*
* GetByID retrieves a user by their ID
*/
func (r *UserRepo) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
