package session

import (
	"testing"
	"time"
)

func TestSessionService_CreateSession(t *testing.T) {
	service := NewSessionService()
	botID := "test-bot-123"

	session := service.CreateSession(botID)

	if session == nil {
		t.Fatal("Expected session to be created, got nil")
	}

	if session.BotID != botID {
		t.Errorf("Expected BotID %s, got %s", botID, session.BotID)
	}

	if session.ID == "" {
		t.Error("Expected session ID to be generated")
	}

	if session.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if session.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}

	// Verify session is stored
	if service.GetSessionCount() != 1 {
		t.Errorf("Expected 1 session, got %d", service.GetSessionCount())
	}
}

func TestSessionService_GetSession(t *testing.T) {
	service := NewSessionService()
	botID := "test-bot-123"

	// Create a session
	created := service.CreateSession(botID)

	// Retrieve the session
	retrieved, err := service.GetSession(created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected session ID %s, got %s", created.ID, retrieved.ID)
	}

	if retrieved.BotID != botID {
		t.Errorf("Expected BotID %s, got %s", botID, retrieved.BotID)
	}
}

func TestSessionService_GetSession_NotFound(t *testing.T) {
	service := NewSessionService()

	_, err := service.GetSession("non-existent-id")
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound, got %v", err)
	}
}

func TestSessionService_UpdateActivity(t *testing.T) {
	service := NewSessionService()
	botID := "test-bot-123"

	session := service.CreateSession(botID)
	originalActivity := session.LastActivityAt

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	err := service.UpdateActivity(session.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Retrieve updated session
	updated, _ := service.GetSession(session.ID)
	if !updated.LastActivityAt.After(originalActivity) {
		t.Error("Expected LastActivityAt to be updated")
	}
}

func TestSessionService_DeleteSession(t *testing.T) {
	service := NewSessionService()
	botID := "test-bot-123"

	session := service.CreateSession(botID)

	// Verify session exists
	if service.GetSessionCount() != 1 {
		t.Errorf("Expected 1 session, got %d", service.GetSessionCount())
	}

	// Delete session
	service.DeleteSession(session.ID)

	// Verify session is deleted
	if service.GetSessionCount() != 0 {
		t.Errorf("Expected 0 sessions, got %d", service.GetSessionCount())
	}

	// Verify session cannot be retrieved
	_, err := service.GetSession(session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound, got %v", err)
	}
}

func TestSessionService_ConcurrentAccess(t *testing.T) {
	service := NewSessionService()
	botID := "test-bot-123"

	// Create multiple sessions concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			service.CreateSession(botID)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all sessions were created
	if service.GetSessionCount() != 10 {
		t.Errorf("Expected 10 sessions, got %d", service.GetSessionCount())
	}
}
