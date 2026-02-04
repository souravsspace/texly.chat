/*
 * This file is auto-generated. Do not edit directly.
 */

/*
 * WidgetConfig holds configuration for the embeddable widget
 */
export interface WidgetConfig {
  theme_color: string;
  initial_message: string;
  position: string;
  bot_avatar: string;
}

/*
 * Bot represents a user's chatbot
 */
export interface Bot {
  id: string;
  user_id: string;
  name: string;
  system_prompt: string;
  allowed_origins: string;
  widget_config: string;
  created_at: string | Date;
  updated_at: string | Date;
  deleted_at: string | Date | null;
}

/*
 * CreateBotRequest holds data for creating a new bot
 */
export interface CreateBotRequest {
  name: string;
  system_prompt: string;
  allowed_origins: string[];
  widget_config: WidgetConfig | null;
}

/*
 * UpdateBotRequest holds data for updating an existing bot
 */
export interface UpdateBotRequest {
  name: string;
  system_prompt: string;
  allowed_origins: string[];
  widget_config: WidgetConfig | null;
}

/*
 * ChatRequest represents an incoming chat message from the user
 */
export interface ChatRequest {
  message: string;
}

/*
 * ChatTokenResponse represents a streaming token or event in SSE format
 */
export interface ChatTokenResponse {
  type: string;
  content: string;
  error: string;
}

/*
 * DocumentChunk represents a chunk of text with its vector embedding
 */
export interface DocumentChunk {
  id: string;
  source_id: string;
  content: string;
  chunk_index: number;
  created_at: string | Date;
  source: Source;
}

/*
 * ChatSession represents an anonymous user session for the widget
 */
export interface ChatSession {
  id: string;
  bot_id: string;
  created_at: string | Date;
  last_activity_at: string | Date;
  expires_at: string | Date;
}

/*
 * CreateSessionRequest holds data for creating a new chat session
 */
export interface CreateSessionRequest {
  bot_id: string;
}

/*
 * SessionResponse represents the response when creating a session
 */
export interface SessionResponse {
  session_id: string;
  bot_id: string;
  expires_at: string | Date;
}

/*
 * SourceType represents the type of data source
 */
export type SourceType = string;

/*
 * SourceStatus represents the processing status of a source
 */
export type SourceStatus = string;

/*
 * Source represents a data source for a bot
 */
export interface Source {
  id: string;
  bot_id: string;
  source_type: SourceType;
  url: string;
  file_path: string;
  original_filename: string;
  content_type: string;
  status: SourceStatus;
  processing_progress: number;
  error_message: string;
  processed_at: string | Date | null;
  created_at: string | Date;
  updated_at: string | Date;
  deleted_at: string | Date | null;
}

/*
 * CreateSourceRequest holds data for creating a new URL source
 */
export interface CreateSourceRequest {
  url: string;
}

/*
 * CreateTextSourceRequest holds data for creating a text source
 */
export interface CreateTextSourceRequest {
  text: string;
  name: string;
}

/*
 * User represents a registered user in the system
 */
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string | Date;
  updated_at: string | Date;
}

/*
 * LoginRequest holds the credentials for user login
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/*
 * SignupRequest holds data for creating a new user
 */
export interface SignupRequest {
  email: string;
  password: string;
  name: string;
}

/*
 * AuthResponse is the response payload for successful authentication
 */
export interface AuthResponse {
  token: string;
  user: User;
}

