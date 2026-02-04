package session

import (
	"errors"
	"sync"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

/*
* SessionService manages anonymous chat sessions in-memory
 */
type SessionService struct {
	sessions map[string]*models.ChatSession
	mu       sync.RWMutex
}

/*
* NewSessionService creates a new session service instance
 */
func NewSessionService() *SessionService {
	service := &SessionService{
		sessions: make(map[string]*models.ChatSession),
	}
	
	// Start background cleanup goroutine
	go service.cleanupExpiredSessions()
	
	return service
}

/*
* CreateSession creates a new anonymous chat session
 */
func (s *SessionService) CreateSession(botID string) *models.ChatSession {
	session := models.NewChatSession(botID)
	
	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()
	
	return session
}

/*
* GetSession retrieves a session by ID
 */
func (s *SessionService) GetSession(sessionID string) (*models.ChatSession, error) {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()
	
	if !exists {
		return nil, ErrSessionNotFound
	}
	
	if session.IsExpired() {
		return nil, ErrSessionExpired
	}
	
	return session, nil
}

/*
* UpdateActivity updates the last activity timestamp for a session
 */
func (s *SessionService) UpdateActivity(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	session, exists := s.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}
	
	if session.IsExpired() {
		return ErrSessionExpired
	}
	
	session.UpdateActivity()
	return nil
}

/*
* DeleteSession removes a session from the store
 */
func (s *SessionService) DeleteSession(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

/*
* cleanupExpiredSessions runs periodically to remove expired sessions
 */
func (s *SessionService) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mu.Lock()
		for id, session := range s.sessions {
			if session.IsExpired() {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

/*
* GetSessionCount returns the current number of active sessions
 */
func (s *SessionService) GetSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}
