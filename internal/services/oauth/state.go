package oauth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type StateService struct {
	redisClient *redis.Client
}

func NewStateService(redisClient *redis.Client) *StateService {
	return &StateService{
		redisClient: redisClient,
	}
}

func (s *StateService) GenerateState() (string, error) {
	state := uuid.New().String()
	// Store with 5 minute expiration
	err := s.redisClient.Set(context.Background(), "oauth_state:"+state, "valid", 5*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return state, nil
}

func (s *StateService) ValidateState(state string) bool {
	key := "oauth_state:" + state
	_, err := s.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		return false
	}
	// Delete state after validation to prevent replay
	s.redisClient.Del(context.Background(), key)
	return true
}
