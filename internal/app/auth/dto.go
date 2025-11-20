package auth

// POST /v1/auth/signup
type SignupRequest struct {
	TempID      string `json:"temp_id,omitempty"`
	Method      string `json:"method" binding:"required"` // "email", "mobile", "google"
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Password    string `json:"password,omitempty"`
	SocialToken string `json:"social_token,omitempty"`
}

type SignupResponse struct {
	UserID              string `json:"user_id"`
	VerificationNeeded  bool   `json:"verification_needed"`
	VerificationChannel string `json:"verification_channel,omitempty"`
}

// POST /v1/auth/verify
type VerifyRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Method string `json:"method" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type VerifyResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

// POST /v1/auth/login
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// POST /v1/auth/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// POST /v1/auth/logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

// POST /v1/auth/resend
type ResendRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Method string `json:"method" binding:"required"`
}

type ResendResponse struct {
	Message string `json:"message"`
}

// UserInfo shared across responses
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
