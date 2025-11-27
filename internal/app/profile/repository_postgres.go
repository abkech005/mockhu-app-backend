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
	// Build dynamic UPDATE query
	query := `UPDATE users SET `
	args := []interface{}{}
	argPos := 1

	for field, value := range updates {
		if argPos > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", field, argPos)
		args = append(args, value)
		argPos++
	}

	// Add WHERE clause and updated_at
	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, userID)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// CheckUsernameExists checks if a username already exists (excluding a specific user)
func (r *PostgresProfileRepository) CheckUsernameExists(ctx context.Context, username string, excludeUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1) AND id != $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, username, excludeUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username: %w", err)
	}

	return exists, nil
}

// UpdateAvatar updates the user's avatar URL
func (r *PostgresProfileRepository) UpdateAvatar(ctx context.Context, userID string, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(ctx, query, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}

// GetPrivacySettings retrieves user's privacy settings
func (r *PostgresProfileRepository) GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error) {
	query := `
		SELECT who_can_message, who_can_see_posts, show_followers_list, show_following_list
		FROM users 
		WHERE id = $1 AND is_active = true
	`

	settings := &PrivacySettings{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&settings.WhoCanMessage,
		&settings.WhoCanSeePosts,
		&settings.ShowFollowersList,
		&settings.ShowFollowingList,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get privacy settings: %w", err)
	}

	return settings, nil
}

// UpdatePrivacySettings updates user's privacy settings
func (r *PostgresProfileRepository) UpdatePrivacySettings(ctx context.Context, userID string, settings *UpdatePrivacyRequest) error {
	// Build dynamic UPDATE query
	updates := make(map[string]interface{})
	
	if settings.WhoCanMessage != "" {
		updates["who_can_message"] = settings.WhoCanMessage
	}
	if settings.WhoCanSeePosts != "" {
		updates["who_can_see_posts"] = settings.WhoCanSeePosts
	}
	if settings.ShowFollowersList != nil {
		updates["show_followers_list"] = *settings.ShowFollowersList
	}
	if settings.ShowFollowingList != nil {
		updates["show_following_list"] = *settings.ShowFollowingList
	}

	// If no updates, return without error
	if len(updates) == 0 {
		return nil
	}

	// Use the existing UpdateProfile method logic
	query := `UPDATE users SET `
	args := []interface{}{}
	argPos := 1

	for field, value := range updates {
		if argPos > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", field, argPos)
		args = append(args, value)
		argPos++
	}

	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, userID)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update privacy settings: %w", err)
	}

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
