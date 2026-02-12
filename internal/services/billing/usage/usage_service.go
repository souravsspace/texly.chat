package usage

import (
	"time"

	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

type UsageService struct {
	db *gorm.DB
}

func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{db: db}
}

// TrackChatMessage records usage for a chat message
func (s *UsageService) TrackChatMessage(userID, botID string) error {
	cost := configs.CalculateMessageCost(1)
	return s.trackUsage(userID, botID, models.UsageTypeChatMessage, 1, cost)
}

// TrackEmbedding records usage for embedding tokens
func (s *UsageService) TrackEmbedding(userID string, tokens int) error {
	cost := configs.CalculateEmbeddingCost(tokens)
	return s.trackUsage(userID, "", models.UsageTypeEmbedding, float64(tokens), cost)
}

// TrackStorage records usage for file storage
func (s *UsageService) TrackStorage(userID string, sizeGB float64) error {
	cost := configs.CalculateStorageCost(sizeGB)
	return s.trackUsage(userID, "", models.UsageTypeStorage, sizeGB, cost)
}

// trackUsage is the internal helper to save the record and update user balance
func (s *UsageService) trackUsage(userID, botID, usageType string, quantity, cost float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create usage record
		record := models.UsageRecord{
			ID:        uuid.New().String(),
			UserID:    userID,
			BotID:     botID,
			Type:      usageType,
			Quantity:  quantity,
			Cost:      cost,
			BilledAt:  time.Time{}, // Zero time means not yet billed/invoiced
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 2. Update user's current period usage
		if err := tx.Model(&models.User{}).Where("id = ?", userID).
			UpdateColumn("current_period_usage", gorm.Expr("current_period_usage + ?", cost)).Error; err != nil {
			return err
		}

		// 3. Deduct from credits if available (Pro tier logic)
		var user models.User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return err
		}

		if user.CreditsBalance > 0 {
			deduction := cost
			if user.CreditsBalance < cost {
				deduction = user.CreditsBalance
			}
			if err := tx.Model(&user).
				UpdateColumn("credits_balance", gorm.Expr("credits_balance - ?", deduction)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetCurrentUsage returns the total cost for the current period
func (s *UsageService) GetCurrentUsage(userID string) (float64, error) {
	var total float64
	// Sum cost of unbilled records
	err := s.db.Model(&models.UsageRecord{}).
		Where("user_id = ? AND billed_at = ?", userID, time.Time{}).
		Select("COALESCE(SUM(cost), 0)").
		Scan(&total).Error
	return total, err
}
