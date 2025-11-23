package onboarding

import "time"

// POST /v1/onboarding/complete - Single comprehensive request for complete onboarding
type CompleteOnboardingRequest struct {
	UserID    string   `json:"user_id" binding:"required"`
	FirstName string   `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string   `json:"last_name" binding:"required,min=2,max=50"`
	Username  string   `json:"username" binding:"required,min=3,max=30"`
	DOB       string   `json:"dob" binding:"required"` // Format: YYYY-MM-DD
	AvatarURL string   `json:"avatar_url,omitempty"`
	Interests []string `json:"interests,omitempty"` // Optional: for future use
}

type CompleteOnboardingResponse struct {
	Success             bool      `json:"success"`
	Message             string    `json:"message"`
	UserID              string    `json:"user_id"`
	ProfileCompleted    bool      `json:"profile_completed"`
	InterestsCount      int       `json:"interests_count"`
	OnboardingCompleted bool      `json:"onboarding_completed"`
	OnboardedAt         time.Time `json:"onboarded_at"`
}

// GET /v1/onboarding/status/:user_id - Check onboarding progress
type OnboardingStatusResponse struct {
	UserID              string `json:"user_id"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneVerified       bool   `json:"phone_verified"`
	ProfileCompleted    bool   `json:"profile_completed"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
	NextStep            string `json:"next_step,omitempty"`
}
