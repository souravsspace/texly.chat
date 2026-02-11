# Phase 6: Authentication Enhancement (OAuth)

## Goal
Add Google and GitHub OAuth 2.0 authentication for improved user onboarding and reduced signup friction.

---

## Backend Tasks

### Step 1: OAuth Dependencies & Configuration

#### 1.1 Add OAuth Libraries
- [ ] Add `golang.org/x/oauth2` to `go.mod`
- [ ] Add `golang.org/x/oauth2/google` for Google OAuth
- [ ] Add `github.com/google/go-github/v58` for GitHub API

#### 1.2 Update Configuration
- [ ] Add OAuth config to `configs/config.go`:
  ```go
  // Google OAuth
  GoogleClientID     string
  GoogleClientSecret string
  GoogleRedirectURL  string // e.g., http://localhost:8080/api/auth/google/callback
  
  // GitHub OAuth
  GitHubClientID     string
  GitHubClientSecret string
  GitHubRedirectURL  string // e.g., http://localhost:8080/api/auth/github/callback
  
  // Frontend URL for redirects
  FrontendURL string // e.g., http://localhost:3000
  ```

#### 1.3 Environment Variables
- [ ] Add to `.env.local`:
  ```
  GOOGLE_CLIENT_ID=your-google-client-id
  GOOGLE_CLIENT_SECRET=your-google-client-secret
  GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
  
  GITHUB_CLIENT_ID=your-github-client-id
  GITHUB_CLIENT_SECRET=your-github-client-secret
  GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback
  
  FRONTEND_URL=http://localhost:3000
  ```

---

### Step 2: Update User Model

#### 2.1 Add OAuth Fields
- [ ] Update `internal/models/user_model.go`:
  ```go
  type User struct {
      ID       string `json:"id" gorm:"primaryKey"`
      Email    string `json:"email" gorm:"uniqueIndex;not null"`
      Password string `json:"-" gorm:"default:null"` // Nullable for OAuth users
      Name     string `json:"name"`
      Avatar   string `json:"avatar"` // Profile picture URL
      
      // OAuth fields
      GoogleID  string `json:"google_id" gorm:"uniqueIndex"`
      GitHubID  int64  `json:"github_id" gorm:"uniqueIndex"`
      AuthProvider string `json:"auth_provider"` // "email", "google", "github"
      
      // Existing fields...
      CreatedAt time.Time      `json:"created_at"`
      UpdatedAt time.Time      `json:"updated_at"`
      DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
  }
  ```
- [ ] Run migration to add new columns
- [ ] Update `BeforeCreate` hook to handle OAuth users (skip password requirement)

---

### Step 3: OAuth Service Layer

#### 3.1 Create OAuth Service
- [ ] Create `internal/services/oauth/oauth_service.go`:
  ```go
  type OAuthService struct {
      googleConfig *oauth2.Config
      githubConfig *oauth2.Config
      userRepo     *user.Repository
      jwtSecret    string
      frontendURL  string
  }
  ```
  - [ ] `GetGoogleAuthURL(state string) string` - Generate auth URL
  - [ ] `GetGitHubAuthURL(state string) string` - Generate auth URL
  - [ ] `HandleGoogleCallback(code string) (*User, string, error)` - Exchange code, create/update user, return JWT
  - [ ] `HandleGitHubCallback(code string) (*User, string, error)` - Exchange code, create/update user, return JWT
  - [ ] `FindOrCreateUser(provider, providerID, email, name, avatar) (*User, error)` - Idempotent user creation

#### 3.2 State Management (CSRF Protection)
- [ ] Create `internal/services/oauth/state.go`:
  - [ ] Generate random state token (UUID)
  - [ ] Store in Redis with 5-minute TTL
  - [ ] Validate state on callback
  - [ ] Delete state after validation

---

### Step 4: OAuth Handlers

#### 4.1 Google OAuth Endpoints
- [ ] Create `internal/handlers/auth/google_handler.go`:
  - [ ] `GET /api/auth/google` - Redirect to Google consent screen
    - [ ] Generate state token
    - [ ] Store in Redis
    - [ ] Redirect to Google OAuth URL
  - [ ] `GET /api/auth/google/callback` - Handle OAuth callback
    - [ ] Validate state token
    - [ ] Exchange code for access token
    - [ ] Fetch user info from Google API
    - [ ] Find or create user in database
    - [ ] Generate JWT token
    - [ ] Redirect to frontend with token in URL fragment

#### 4.2 GitHub OAuth Endpoints
- [ ] Create `internal/handlers/auth/github_handler.go`:
  - [ ] `GET /api/auth/github` - Redirect to GitHub consent screen
    - [ ] Generate state token
    - [ ] Store in Redis
    - [ ] Redirect to GitHub OAuth URL
  - [ ] `GET /api/auth/github/callback` - Handle OAuth callback
    - [ ] Validate state token
    - [ ] Exchange code for access token
    - [ ] Fetch user info from GitHub API
    - [ ] Fetch primary email if not public
    - [ ] Find or create user in database
    - [ ] Generate JWT token
    - [ ] Redirect to frontend with token in URL fragment

#### 4.3 Register Routes
- [ ] Update `internal/server/server.go`:
  ```go
  auth := router.Group("/api/auth")
  {
      // Existing email/password routes
      auth.POST("/signup", authHandler.SignUp)
      auth.POST("/login", authHandler.Login)
      
      // New OAuth routes
      auth.GET("/google", authHandler.GoogleAuth)
      auth.GET("/google/callback", authHandler.GoogleCallback)
      auth.GET("/github", authHandler.GitHubAuth)
      auth.GET("/github/callback", authHandler.GitHubCallback)
  }
  ```

---

### Step 5: Account Linking (Optional Future Enhancement)

#### 5.1 Link Existing Accounts
- [ ] Allow users with email/password to link Google/GitHub
- [ ] Prevent duplicate accounts with same email
- [ ] Merge logic:
  - [ ] If email exists, update `google_id` or `github_id`
  - [ ] If new email, create new user
- [ ] Add endpoints:
  - [ ] `POST /api/user/link/google` - Link Google account
  - [ ] `POST /api/user/link/github` - Link GitHub account
  - [ ] `DELETE /api/user/unlink/:provider` - Unlink OAuth provider

---

## Frontend Tasks

### Step 1: Update Auth UI Components

#### 1.1 Social Login Buttons
- [ ] Create `ui/src/components/auth/social-login-buttons.tsx`:
  - [ ] Google sign-in button with Google branding
  - [ ] GitHub sign-in button with GitHub branding
  - [ ] Use shadcn/ui Button component
  - [ ] Add loading states
  - [ ] Handle errors

#### 1.2 Update Login Page
- [ ] Update `ui/src/routes/_auth/login.tsx`:
  - [ ] Add social login buttons above email/password form
  - [ ] Add "OR" divider between social and email login
  - [ ] Style consistently with existing design
  - [ ] Add error toast for OAuth failures

#### 1.3 Update Signup Page
- [ ] Update `ui/src/routes/_auth/signup.tsx`:
  - [ ] Add social signup buttons above email/password form
  - [ ] Add "OR" divider
  - [ ] Update copy: "Sign up with Google" / "Sign up with GitHub"

---

### Step 2: OAuth Flow Implementation

#### 2.1 Handle OAuth Redirects
- [ ] Create `ui/src/routes/_auth/oauth-callback.tsx`:
  - [ ] Parse JWT token from URL fragment
  - [ ] Store token in localStorage/secure storage
  - [ ] Fetch user profile
  - [ ] Redirect to dashboard
  - [ ] Show error if token invalid

#### 2.2 API Client Updates
- [ ] Update `ui/src/api/auth.ts`:
  ```ts
  export const initiateGoogleAuth = () => {
    window.location.href = `${API_URL}/api/auth/google`;
  };
  
  export const initiateGitHubAuth = () => {
    window.location.href = `${API_URL}/api/auth/github`;
  };
  ```

---

### Step 3: User Profile Updates

#### 3.1 Display OAuth Info
- [ ] Update `ui/src/routes/dashboard/settings.tsx`:
  - [ ] Show connected OAuth providers
  - [ ] Display avatar from OAuth provider
  - [ ] Show "Connected via Google" badge if applicable
  - [ ] Add "Link Google Account" / "Link GitHub Account" buttons
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

#### GitHub OAuth Setup
- [ ] Go to GitHub Settings > Developer Settings > OAuth Apps
- [ ] Create new OAuth app
- [ ] Set Authorization callback URL:
  - [ ] `http://localhost:8080/api/auth/github/callback` (dev)
  - [ ] `https://yourdomain.com/api/auth/github/callback` (prod)
- [ ] Copy Client ID and generate Client Secret
- [ ] Add to `.env.local`

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
- [ ] Mock GitHub OAuth flow
- [ ] Test callback error handling
- [ ] Test state mismatch scenarios

### Manual Testing
- [ ] Sign up with Google
- [ ] Sign up with GitHub
- [ ] Login with existing Google account
- [ ] Login with existing GitHub account
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
