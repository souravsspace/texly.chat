package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/cache"
	"gorm.io/gorm"
)

/*
* UserRepo handles database operations for users
 */
type UserRepo struct {
	db    *gorm.DB
	cache *cache.CacheService
}

/*
* NewUserRepo creates a new UserRepo instance
 */
func NewUserRepo(db *gorm.DB, cache *cache.CacheService) *UserRepo {
	return &UserRepo{db: db, cache: cache}
}

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
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.UserEmailCacheKey, email)

	// Try cache first
	var user models.User
	if err := r.cache.GetJSON(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// Cache miss - query database
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, user, cache.UserCacheTTL)
	return &user, nil
}

/*
* GetByID retrieves a user by their ID
 */
func (r *UserRepo) GetByID(id string) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.UserCacheKey, id)

	// Try cache first
	var user models.User
	if err := r.cache.GetJSON(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// Cache miss - query database
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, user, cache.UserCacheTTL)
	return &user, err
}
