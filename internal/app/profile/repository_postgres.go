package profile

import (
	"context"
	"fmt"

	"mockhu-app-backend/internal/app/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresProfileRepository implements ProfileRepository using PostgreSQL
type PostgresProfileRepository struct {
	db *pgxpool.Pool
}

// NewPostgresProfileRepository creates a new PostgreSQL profile repository
func NewPostgresProfileRepository(db *pgxpool.Pool) ProfileRepository {
	return &PostgresProfileRepository{db: db}
}

// GetProfileByID retrieves a user's profile by ID
func (r *PostgresProfileRepository) GetProfileByID(ctx context.Context, userID string) (*auth.User, error) {
	query := `
		SELECT 
			id, 
			COALESCE(email, ''), 
			email_verified, 
			COALESCE(phone, ''), 
			phone_verified,
			COALESCE(username, ''), 
			COALESCE(first_name, ''), 
			COALESCE(last_name, ''), 
			dob, 
			COALESCE(avatar_url, ''),
			COALESCE(bio, ''), 
			institution_id,
			who_can_message, 
			who_can_see_posts, 
			show_followers_list, 
			show_following_list,
			is_active, 
			onboarding_completed, 
			onboarded_at,
			last_login_at, 
			created_at, 
			updated_at
		FROM users 
		WHERE id = $1 AND is_active = true
	`

	user := &auth.User{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.EmailVerified,
		&user.Phone,
		&user.PhoneVerified,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&user.AvatarURL,
		&user.Bio,
		&user.InstitutionID,
		&user.WhoCanMessage,
		&user.WhoCanSeePosts,
		&user.ShowFollowersList,
		&user.ShowFollowingList,
		&user.IsActive,
		&user.OnboardingCompleted,
		&user.OnboardedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}

// UpdateProfile updates user profile fields
func (r *PostgresProfileRepository) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	// TODO: Implement in Phase 4
	return nil
}

// CheckUsernameExists checks if a username already exists (excluding a specific user)
func (r *PostgresProfileRepository) CheckUsernameExists(ctx context.Context, username string, excludeUserID string) (bool, error) {
	// TODO: Implement in Phase 4
	return false, nil
}

// UpdateAvatar updates the user's avatar URL
func (r *PostgresProfileRepository) UpdateAvatar(ctx context.Context, userID string, avatarURL string) error {
	// TODO: Implement in Phase 5
	return nil
}

// GetPrivacySettings retrieves user's privacy settings
func (r *PostgresProfileRepository) GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error) {
	// TODO: Implement in Phase 6
	return nil, nil
}

// UpdatePrivacySettings updates user's privacy settings
func (r *PostgresProfileRepository) UpdatePrivacySettings(ctx context.Context, userID string, settings *UpdatePrivacyRequest) error {
	// TODO: Implement in Phase 6
	return nil
}

// GetMutualConnections retrieves mutual connections between two users
func (r *PostgresProfileRepository) GetMutualConnections(ctx context.Context, user1ID, user2ID string, limit, offset int) ([]*auth.User, error) {
	// TODO: Implement in Phase 7
	return nil, nil
}

// GetMutualConnectionsCount retrieves the count of mutual connections
func (r *PostgresProfileRepository) GetMutualConnectionsCount(ctx context.Context, user1ID, user2ID string) (int, error) {
	// TODO: Implement in Phase 7
	return 0, nil
}
