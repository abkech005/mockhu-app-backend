package auth

import "time"

// User represents the domain model
type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	Phone         string     `json:"phone"`
	PhoneVerified bool       `json:"phone_verified"`
	Username      string     `json:"username"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	DOB           time.Time  `json:"dob"`
	PasswordHash  string     `json:"-"` // Never expose in JSON
	AvatarURL     string     `json:"avatar_url"`
	IsActive      bool       `json:"is_active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
