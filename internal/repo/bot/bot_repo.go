package bot

import (
	"errors"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
* BotRepo handles database operations for bots
 */
type BotRepo struct {
	db *gorm.DB
}

/*
* NewBotRepo creates a new BotRepo instance
 */
func NewBotRepo(db *gorm.DB) *BotRepo {
	return &BotRepo{db: db}
}

/*
* Create inserts a new bot into the database
 */
func (r *BotRepo) Create(bot *models.Bot) error {
	return r.db.Create(bot).Error
}

/*
* GetByUserID retrieves all bots for a specific user
 */
func (r *BotRepo) GetByUserID(userID string) ([]models.Bot, error) {
	var bots []models.Bot
	err := r.db.Where("user_id = ?", userID).Find(&bots).Error
	return bots, err
}

/*
* GetByID retrieves a bot by its ID and UserID
 */
func (r *BotRepo) GetByID(id string, userID string) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&bot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bot, nil
}

/*
* Update updates an existing bot
 */
func (r *BotRepo) Update(bot *models.Bot) error {
	return r.db.Save(bot).Error
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
	return nil
}

/*
* GetByIDPublic retrieves a bot by its ID without user authentication
* Used for public widget access and CORS validation
 */
func (r *BotRepo) GetByIDPublic(id string) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("id = ?", id).First(&bot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bot, nil
}
