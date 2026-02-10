package analytics

import (
	"context"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	messageRepo "github.com/souravsspace/texly.chat/internal/repo/message"
)

/*
 * AnalyticsService provides analytics aggregation and reporting
 */
type AnalyticsService struct {
	messageRepo *messageRepo.MessageRepository
}

/*
 * NewAnalyticsService creates a new analytics service instance
 */
func NewAnalyticsService(messageRepo *messageRepo.MessageRepository) *AnalyticsService {
	return &AnalyticsService{
		messageRepo: messageRepo,
	}
}

/*
 * GetBotAnalytics retrieves comprehensive analytics for a specific bot
 */
func (s *AnalyticsService) GetBotAnalytics(ctx context.Context, botID string) (*models.BotAnalytics, error) {
	return s.messageRepo.GetBotAnalytics(ctx, botID)
}

/*
 * GetBotDailyStats retrieves daily message statistics for a bot
 */
func (s *AnalyticsService) GetBotDailyStats(ctx context.Context, botID string, days int) ([]models.MessageStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	return s.messageRepo.GetDailyStats(ctx, botID, startDate, endDate)
}

/*
 * GetUserAnalytics retrieves analytics for all bots owned by a user
 */
func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, userID string) (map[string]*models.BotAnalytics, error) {
	return s.messageRepo.GetUserAnalytics(ctx, userID)
}

/*
 * GetSessionMessages retrieves all messages for a specific session
 */
func (s *AnalyticsService) GetSessionMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	return s.messageRepo.GetBySessionID(ctx, sessionID)
}
