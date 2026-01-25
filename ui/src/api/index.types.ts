/*
 * This file is auto-generated. Do not edit directly.
 */

/*
 * Post represents a user-created post
 */
export interface Post {
  id: string;
  user_id: string;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
}

/*
 * CreatePostRequest holds data for creating a new post
 */
export interface CreatePostRequest {
  title: string;
  content: string;
}

/*
 * UpdatePostRequest holds data for updating an existing post
 */
export interface UpdatePostRequest {
  title: string;
  content: string;
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
