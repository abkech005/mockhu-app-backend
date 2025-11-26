package follow

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFollowRepository implements FollowRepository for PostgreSQL
type PostgresFollowRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresFollowRepository creates a new PostgreSQL follow repository
func NewPostgresFollowRepository(pool *pgxpool.Pool) FollowRepository {
	return &PostgresFollowRepository{pool: pool}
}

// Follow creates a follow relationship
func (r *PostgresFollowRepository) Follow(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO user_follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, followerID, followingID)
	return err
}

// Unfollow removes a follow relationship
func (r *PostgresFollowRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND following_id = $2
	`

	_, err := r.pool.Exec(ctx, query, followerID, followingID)
	return err
}

// IsFollowing checks if follower follows following
func (r *PostgresFollowRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

// GetFollowers returns users who follow the given user
func (r *PostgresFollowRepository) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at
		FROM user_follows
		WHERE following_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		f := &Follow{}
		err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, f)
	}

	return follows, rows.Err()
}

// GetFollowing returns users that the given user follows
func (r *PostgresFollowRepository) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at
		FROM user_follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		f := &Follow{}
		err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, f)
	}

	return follows, rows.Err()
}

// GetFollowerCount returns the number of followers
func (r *PostgresFollowRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM user_follows WHERE following_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// GetFollowingCount returns the number of users being followed
func (r *PostgresFollowRepository) GetFollowingCount(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// GetFollowStats returns both follower and following counts
func (r *PostgresFollowRepository) GetFollowStats(ctx context.Context, userID string) (*FollowStats, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM user_follows WHERE following_id = $1) as follower_count,
			(SELECT COUNT(*) FROM user_follows WHERE follower_id = $1) as following_count
	`

	stats := &FollowStats{}
	err := r.pool.QueryRow(ctx, query, userID).Scan(&stats.FollowerCount, &stats.FollowingCount)
	return stats, err
}
