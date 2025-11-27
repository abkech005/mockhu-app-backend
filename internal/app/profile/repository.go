package profile

import (
	"context"

	"mockhu-app-backend/internal/app/auth"
)

// ProfileRepository defines the interface for profile data operations
type ProfileRepository interface {
	// Profile operations
	GetProfileByID(ctx context.Context, userID string) (*auth.User, error)
	UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) error
	CheckUsernameExists(ctx context.Context, username string, excludeUserID string) (bool, error)

	// Avatar operations
	UpdateAvatar(ctx context.Context, userID string, avatarURL string) error

	// Privacy operations
	GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, userID string, settings *UpdatePrivacyRequest) error

	// Mutual connections
	GetMutualConnections(ctx context.Context, user1ID, user2ID string, limit, offset int) ([]*auth.User, error)
	GetMutualConnectionsCount(ctx context.Context, user1ID, user2ID string) (int, error)
}
