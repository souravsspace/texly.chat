package credits

import (
	"fmt"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

type CreditsService struct {
	db *gorm.DB
}

func NewCreditsService(db *gorm.DB) *CreditsService {
	return &CreditsService{db: db}
}

// AddCredits adds credits to a user's balance
func (s *CreditsService) AddCredits(userID string, amount float64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("credits_balance", gorm.Expr("credits_balance + ?", amount)).Error
}

// DeductCredits attempts to deduct credits from a user's balance.
// Returns nil if successful, or error if insufficient funds or DB error.
func (s *CreditsService) DeductCredits(userID string, amount float64) error {
	// Use a transaction to ensure atomicity
	return s.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Select("credits_balance").First(&user, "id = ?", userID).Error; err != nil {
			return err
		}

		if user.CreditsBalance < amount {
			return fmt.Errorf("insufficient credits: balance %.2f, required %.2f", user.CreditsBalance, amount)
		}

		return tx.Model(&models.User{}).Where("id = ?", userID).
			UpdateColumn("credits_balance", gorm.Expr("credits_balance - ?", amount)).Error
	})
}

// RefreshMonthlyCredits resets credits for Pro users to the monthly allocation
// Note: Unused credits do NOT roll over (per configs/pricing.go)
func (s *CreditsService) RefreshMonthlyCredits(userID string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("credits_balance", configs.ProIncludedCredits).Error // Set to exactly $20
}

// GetCreditsBalance returns the current credit balance
func (s *CreditsService) GetCreditsBalance(userID string) (float64, error) {
	var user models.User
	if err := s.db.Select("credits_balance").First(&user, "id = ?", userID).Error; err != nil {
		return 0, err
	}
	return user.CreditsBalance, nil
}
