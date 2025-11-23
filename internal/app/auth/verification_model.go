package auth

import "time"

// VerificationCode represents a verification code for email or phone verification
type VerificationCode struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Contact   string     `json:"contact"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// Constants for verification types
const (
	VerificationTypeEmail = "email"
	VerificationTypePhone = "phone"
)
