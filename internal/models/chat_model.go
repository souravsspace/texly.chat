package models

/*
 * ChatRequest represents an incoming chat message from the user
 */
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

/*
 * ChatTokenResponse represents a streaming token or event in SSE format
 */
type ChatTokenResponse struct {
	Type    string `json:"type"`              // "token" | "done" | "error"
	Content string `json:"content,omitempty"` // Token content for type="token"
	Error   string `json:"error,omitempty"`   // Error message for type="error"
}
