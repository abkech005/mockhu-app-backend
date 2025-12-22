package auth

import "context"

// EmailService defines the interface for sending emails
type EmailService interface {
	SendVerificationEmail(ctx context.Context, toEmail, code string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, code string) error
}
