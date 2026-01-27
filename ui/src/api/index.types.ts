/*
 * This file is auto-generated. Do not edit directly.
 */

/*
 * Bot represents a user's chatbot
 */
export interface Bot {
  id: string;
  user_id: string;
  name: string;
  system_prompt: string;
  created_at: string;
  updated_at: string;
  deleted_at: any;
}

/*
 * CreateBotRequest holds data for creating a new bot
 */
export interface CreateBotRequest {
  name: string;
  system_prompt: string;
}

/*
 * UpdateBotRequest holds data for updating an existing bot
 */
export interface UpdateBotRequest {
  name: string;
  system_prompt: string;
}

/*
 * DocumentChunk represents a chunk of text with its vector embedding
 */
export interface DocumentChunk {
  id: string;
  source_id: string;
  content: string;
  chunk_index: number;
  created_at: string;
}

/*
 * User represents a registered user in the system
 */
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
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

