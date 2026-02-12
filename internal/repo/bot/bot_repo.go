package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/cache"
	"gorm.io/gorm"
)

/*
* BotRepo handles database operations for bots
 */
type BotRepo struct {
	db    *gorm.DB
	cache *cache.CacheService
}

/*
* NewBotRepo creates a new BotRepo instance
 */
func NewBotRepo(db *gorm.DB, cache *cache.CacheService) *BotRepo {
	return &BotRepo{db: db, cache: cache}
}

/*
* Create inserts a new bot into the database
 */
func (r *BotRepo) Create(bot *models.Bot) error {
	if err := r.db.Create(bot).Error; err != nil {
		return err
	}
	// Invalidate bot list cache for this user
	_ = r.cache.DeletePattern(context.Background(), fmt.Sprintf(cache.BotListCacheKey, bot.UserID))
	return nil
}

/*
* GetByUserID retrieves all bots for a specific user
 */
func (r *BotRepo) GetByUserID(userID string) ([]models.Bot, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.BotListCacheKey, userID)

	// Try cache first
	var bots []models.Bot
	if err := r.cache.GetJSON(ctx, cacheKey, &bots); err == nil {
		return bots, nil
	}

	// Cache miss - query database
	err := r.db.Where("user_id = ?", userID).Find(&bots).Error
	if err != nil {
		return bots, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, bots, cache.BotListCacheTTL)
	return bots, nil
}

/*
* GetByID retrieves a bot by its ID and UserID
 */
func (r *BotRepo) GetByID(id string, userID string) (*models.Bot, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.BotCacheKey, id)

	// Try cache first
	var bot models.Bot
	if err := r.cache.GetJSON(ctx, cacheKey, &bot); err == nil {
		// Verify user ownership from cached data
		if bot.UserID == userID {
			return &bot, nil
		}
		// Cache hit but different user, fall through to DB
	}

	// Cache miss - query database
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&bot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, bot, cache.BotCacheTTL)
	return &bot, nil
}

/*
* Update updates an existing bot
 */
func (r *BotRepo) Update(bot *models.Bot) error {
	if err := r.db.Save(bot).Error; err != nil {
		return err
	}
	// Invalidate caches
	ctx := context.Background()
	_ = r.cache.Delete(ctx, fmt.Sprintf(cache.BotCacheKey, bot.ID))
	_ = r.cache.DeletePattern(ctx, fmt.Sprintf(cache.BotListCacheKey, bot.UserID))
	return nil
}

/*
* Delete removes a bot from the database
 */
func (r *BotRepo) Delete(id string, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Bot{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	// Invalidate caches
	ctx := context.Background()
	_ = r.cache.Delete(ctx, fmt.Sprintf(cache.BotCacheKey, id))
	_ = r.cache.DeletePattern(ctx, fmt.Sprintf(cache.BotListCacheKey, userID))
	return nil
}

/*
* GetByIDPublic retrieves a bot by its ID without user authentication
* Used for public widget access and CORS validation
 */
func (r *BotRepo) GetByIDPublic(id string) (*models.Bot, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cache.BotCacheKey, id)

	// Try cache first
	var bot models.Bot
	if err := r.cache.GetJSON(ctx, cacheKey, &bot); err == nil {
		return &bot, nil
	}

	// Cache miss - query database
	err := r.db.Where("id = ?", id).First(&bot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Populate cache
	_ = r.cache.SetJSON(ctx, cacheKey, bot, cache.BotCacheTTL)
	return &bot, nil
}
