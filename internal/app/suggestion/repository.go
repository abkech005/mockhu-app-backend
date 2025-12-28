package suggestion

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for suggestions
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new suggestion repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetSuggestedUsers returns users with shared interests, excluding already followed users
func (r *Repository) GetSuggestedUsers(ctx context.Context, userID string, limit int) ([]SuggestedUser, error) {
	// This query:
	// 1. Finds users who share interests with the requesting user
	// 2. Counts shared interests
	// 3. Gets follower count
	// 4. Excludes already followed users
	// 5. Excludes the requesting user
	// 6. Orders by shared interest count (highest first)
	query := `
		WITH my_interests AS (
			SELECT interest_id FROM user_interests WHERE user_id = $1
		),
		user_follower_counts AS (
			SELECT following_id, COUNT(*) as follower_count
			FROM user_follows
			GROUP BY following_id
		),
		shared_interest_users AS (
			SELECT 
				u.id,
				u.username,
				u.first_name,
				u.last_name,
				COALESCE(u.avatar_url, '') as avatar_url,
				COALESCE(u.bio, '') as bio,
				COALESCE(u.place, '') as place,
				COUNT(DISTINCT other_ui.interest_id) as shared_count,
				ARRAY_AGG(DISTINCT i.slug) as shared_slugs,
				COALESCE(ufc.follower_count, 0) as follower_count
			FROM users u
			INNER JOIN user_interests other_ui ON other_ui.user_id = u.id
			INNER JOIN my_interests mi ON other_ui.interest_id = mi.interest_id
			INNER JOIN interests i ON i.id = other_ui.interest_id
			LEFT JOIN user_follower_counts ufc ON ufc.following_id = u.id
			WHERE u.id != $1
				AND u.is_active = true
				AND u.onboarding_completed = true
				AND NOT EXISTS (
					SELECT 1 FROM user_follows 
					WHERE follower_id = $1 AND following_id = u.id
				)
			GROUP BY u.id, u.username, u.first_name, u.last_name, u.avatar_url, u.bio, u.place, ufc.follower_count
		)
		SELECT 
			id, username, first_name, last_name, avatar_url, bio, place,
			shared_count, shared_slugs, follower_count
		FROM shared_interest_users
		ORDER BY shared_count DESC, follower_count DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested users: %w", err)
	}
	defer rows.Close()

	var suggestions []SuggestedUser
	for rows.Next() {
		var s SuggestedUser
		var sharedSlugs []string

		err := rows.Scan(
			&s.UserID, &s.Username, &s.FirstName, &s.LastName,
			&s.AvatarURL, &s.Bio, &s.Place,
			&s.SharedInterestsCount, &sharedSlugs, &s.FollowerCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan suggested user: %w", err)
		}

		s.SharedInterests = sharedSlugs
		s.Reason = buildReason(s.SharedInterestsCount, s.Place)
		suggestions = append(suggestions, s)
	}

	return suggestions, nil
}

// buildReason creates a human-readable reason for the suggestion
func buildReason(sharedCount int, place string) string {
	parts := []string{}

	if sharedCount > 0 {
		if sharedCount == 1 {
			parts = append(parts, "shares 1 interest with you")
		} else {
			parts = append(parts, fmt.Sprintf("shares %d interests with you", sharedCount))
		}
	}

	if place != "" {
		parts = append(parts, fmt.Sprintf("from %s", place))
	}

	if len(parts) == 0 {
		return "suggested for you"
	}

	return strings.Join(parts, " • ")
}
