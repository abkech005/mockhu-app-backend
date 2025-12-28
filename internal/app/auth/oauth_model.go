package auth

import "time"

// OAuthAccount represents a linked OAuth provider account
type OAuthAccount struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	ProviderEmail  string     `json:"provider_email"`
	AccessToken    string     `json:"-"`
	RefreshToken   string     `json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
