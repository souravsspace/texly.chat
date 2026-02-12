# Phase 6: Authentication Enhancement (OAuth) ✅

**Status**: Completed  
**Completed**: February 12, 2025

## Goal
Add Google OAuth 2.0 authentication for improved user onboarding and reduced signup friction.

---

## Backend Tasks

### Step 1: OAuth Dependencies & Configuration

#### 1.1 Add OAuth Libraries
- [x] Add `golang.org/x/oauth2` to `go.mod`
- [x] Add `golang.org/x/oauth2/google` for Google OAuth

#### 1.2 Update Configuration
- [x] Add OAuth config to `configs/config.go`:
  ```go
  // Google OAuth
  GoogleClientID     string
  GoogleClientSecret string
  GoogleRedirectURL  string // e.g., http://localhost:8080/api/auth/google/callback

  // Frontend URL for redirects
  FrontendURL string // e.g., http://localhost:3000
  ```

#### 1.3 Environment Variables
- [x] Add to `.env.local`:
  ```
  GOOGLE_CLIENT_ID=your-google-client-id
  GOOGLE_CLIENT_SECRET=your-google-client-secret
  GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

  FRONTEND_URL=http://localhost:3000
  ```

---

### Step 2: Update User Model

#### 2.1 Add OAuth Fields
- [x] Update `internal/models/user_model.go`:
  ```go
  type User struct {
      ID       string `json:"id" gorm:"primaryKey"`
      Email    string `json:"email" gorm:"uniqueIndex;not null"`
      Password string `json:"-" gorm:"default:null"` // Nullable for OAuth users
      Name     string `json:"name"`
      Avatar   string `json:"avatar"` // Profile picture URL
      
      // OAuth fields
      GoogleID  string `json:"google_id" gorm:"uniqueIndex"`
      AuthProvider string `json:"auth_provider"` // "email", "google"
      
      // Existing fields...
      CreatedAt time.Time      `json:"created_at"`
      UpdatedAt time.Time      `json:"updated_at"`
      DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
  }
  ```
- [x] Run migration to add new columns
- [x] Update `BeforeCreate` hook to handle OAuth users (skip password requirement)

---

### Step 3: OAuth Service Layer

#### 3.1 Create OAuth Service
- [x] Create `internal/services/oauth/oauth_service.go`:
  ```go
  type OAuthService struct {
      googleConfig *oauth2.Config
      userRepo     *user.Repository
      jwtSecret    string
      frontendURL  string
  }
  ```
  - [x] `GetGoogleAuthURL(state string) string` - Generate auth URL
  - [x] `FindOrCreateUser(provider, providerID, email, name, avatar) (*User, error)` - Idempotent user creation

#### 3.2 State Management (CSRF Protection)
- [x] Create `internal/services/oauth/state.go`:
  - [x] Generate random state token (UUID)
  - [x] Store in Redis with 5-minute TTL
  - [x] Validate state on callback
  - [x] Delete state after validation

---

### Step 4: OAuth Handlers

#### 4.1 Google OAuth Endpoints
- [x] Create `internal/handlers/auth/google_handler.go`:
  - [x] `GET /api/auth/google` - Redirect to Google consent screen
    - [x] Generate state token
    - [x] Store in Redis
    - [x] Redirect to Google OAuth URL
  - [x] `GET /api/auth/google/callback` - Handle OAuth callback
    - [x] Validate state token
    - [x] Exchange code for access token
    - [x] Fetch user info from Google API
    - [x] Find or create user in database
    - [x] Generate JWT token
    - [x] Redirect to frontend with token in URL fragment

#### 4.3 Register Routes
- [x] Update `internal/server/server.go`:
  ```go
  auth := router.Group("/api/auth")
  {
      // Existing email/password routes
      auth.POST("/signup", authHandler.SignUp)
      auth.POST("/login", authHandler.Login)
      
      // New OAuth routes
      auth.GET("/google", authHandler.GoogleAuth)
      auth.GET("/google/callback", authHandler.GoogleCallback)
  }
  ```

---

### Step 5: Account Linking (Optional Future Enhancement)

#### 5.1 Link Existing Accounts
- [ ] Allow users with email/password to link Google
- [ ] Prevent duplicate accounts with same email
- [ ] Merge logic:
  - [ ] If email exists, update `google_id`
  - [ ] If new email, create new user
- [ ] Add endpoints:
  - [ ] `POST /api/user/link/google` - Link Google account
  - [ ] `DELETE /api/user/unlink/:provider` - Unlink OAuth provider

---

## Frontend Tasks

### Step 1: Update Auth UI Components

#### 1.1 Social Login Buttons
- [x] Create `ui/src/components/auth/social-login-buttons.tsx`:
  - [x] Google sign-in button with Google branding
  - [x] Use shadcn/ui Button component
  - [x] Add loading states
  - [x] Handle errors

#### 1.2 Update Login Page
- [x] Update `ui/src/routes/_auth/login.tsx`:
  - [x] Add social login buttons above email/password form
  - [x] Add "OR" divider between social and email login
  - [x] Style consistently with existing design
  - [x] Add error toast for OAuth failures

#### 1.3 Update Signup Page
- [x] Update `ui/src/routes/_auth/signup.tsx`:
  - [x] Add social signup buttons above email/password form
  - [x] Add "OR" divider
  - [x] Update copy: "Sign up with Google"

---

### Step 2: OAuth Flow Implementation

#### 2.1 Handle OAuth Redirects
- [x] Create `ui/src/routes/_auth/oauth-callback.tsx`:
  - [x] Parse JWT token from URL fragment
  - [x] Store token in localStorage/secure storage
  - [x] Fetch user profile
  - [x] Redirect to dashboard
  - [x] Show error if token invalid

#### 2.2 API Client Updates
- [x] Update `ui/src/api/auth.ts`:
  ```ts
  export const initiateGoogleAuth = () => {
    window.location.href = `${API_URL}/api/auth/google`;
  };
  ```

---

### Step 3: User Profile Updates

#### 3.1 Display OAuth Info
- [ ] Update `ui/src/routes/dashboard/settings.tsx`:
  - [ ] Show connected OAuth providers
  - [ ] Display avatar from OAuth provider
  - [ ] Show "Connected via Google" badge if applicable
  - [ ] Add "Link Google Account" button
  - [ ] Add "Unlink" option for connected providers

#### 3.2 Avatar Display
- [ ] Update `ui/src/components/layout/header.tsx`:
  - [ ] Display user avatar from OAuth provider
  - [ ] Fallback to initials if no avatar

---

## Security Considerations

### 3rd Party OAuth Setup

#### Google OAuth Setup
- [ ] Create project in Google Cloud Console
- [ ] Enable Google+ API
- [ ] Create OAuth 2.0 credentials
- [ ] Add authorized redirect URIs:
  - [ ] `http://localhost:8080/api/auth/google/callback` (dev)
  - [ ] `https://yourdomain.com/api/auth/google/callback` (prod)
- [ ] Copy Client ID and Client Secret to `.env.local`

---

### Security Checklist
- [ ] State tokens prevent CSRF attacks
- [ ] State stored in Redis with 5-minute expiration
- [ ] Validate state on every callback
- [ ] Use HTTPS in production for OAuth redirects
- [ ] Don't expose OAuth secrets in frontend
- [ ] Validate email domains if needed (e.g., enterprise SSO)
- [ ] Rate limit OAuth endpoints to prevent abuse

---

## Testing Tasks

### Unit Tests
- [ ] Test OAuth service user creation
- [ ] Test state token generation and validation
- [ ] Test duplicate email handling
- [ ] Test account linking logic

### Integration Tests
- [ ] Mock Google OAuth flow
- [ ] Test callback error handling
- [ ] Test state mismatch scenarios

### Manual Testing
- [ ] Sign up with Google
- [ ] Login with existing Google account
- [ ] Try to create duplicate account with same email
- [ ] Test account linking (if implemented)

---

## Success Metrics

- [ ] OAuth signup success rate >95%
- [ ] Reduced signup time from 2 minutes to <30 seconds
- [ ] >50% of new users choose social login
- [ ] Zero CSRF attack vulnerabilities
- [ ] No duplicate accounts with same email

---

## Documentation Updates

- [ ] Update README with OAuth setup instructions
- [ ] Document environment variables for OAuth
- [ ] Add troubleshooting guide for common OAuth errors
- [ ] Update API documentation with new endpoints

---

**Created**: February 11, 2025
