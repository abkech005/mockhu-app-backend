package auth

import (
	"context"
)

// VerificationRepository defines methods for verification code data access
type VerificationRepository interface {
	// Create inserts a new verification code
	Create(ctx context.Context, verification *VerificationCode) error

	// FindByCodeAndType finds an active (unused, non-expired) verification code
	FindByCodeAndType(ctx context.Context, code string, verificationType string) (*VerificationCode, error)

	// FindActiveByContactAndType finds the latest active code for a contact
	FindActiveByContactAndType(ctx context.Context, contact string, verificationType string) (*VerificationCode, error)

	// MarkAsUsed marks a verification code as used
	MarkAsUsed(ctx context.Context, id string) error

	// InvalidatePreviousCodes marks all previous codes for a contact as used (when generating new code)
	InvalidatePreviousCodes(ctx context.Context, userID string, verificationType string) error

	// CleanupExpired deletes expired verification codes (for maintenance)
	CleanupExpired(ctx context.Context) (int64, error)
}
