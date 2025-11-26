package share

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresShareRepository implements ShareRepository for PostgreSQL
type PostgresShareRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresShareRepository creates a new PostgreSQL share repository
func NewPostgresShareRepository(pool *pgxpool.Pool) ShareRepository {
	return &PostgresShareRepository{pool: pool}
}

// Create inserts a new share into the database
func (r *PostgresShareRepository) Create(ctx context.Context, share *Share) error {
	query := `
		INSERT INTO post_shares (post_id, user_id, shared_to_type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		share.PostID,
		share.UserID,
		share.SharedToType,
	).Scan(&share.ID, &share.CreatedAt)

	return err
}

// GetByID retrieves a single share by its ID
func (r *PostgresShareRepository) GetByID(ctx context.Context, id string) (*Share, error) {
	query := `
		SELECT id, post_id, user_id, shared_to_type, created_at
		FROM post_shares
		WHERE id = $1
	`

	share := &Share{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&share.ID,
		&share.PostID,
		&share.UserID,
		&share.SharedToType,
		&share.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return share, nil
}

// GetByPostID retrieves all shares for a post
func (r *PostgresShareRepository) GetByPostID(ctx context.Context, postID string, limit, offset int) ([]*Share, error) {
	query := `
		SELECT id, post_id, user_id, shared_to_type, created_at
		FROM post_shares
		WHERE post_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*Share

	for rows.Next() {
		share := &Share{}
		err := rows.Scan(
			&share.ID,
			&share.PostID,
			&share.UserID,
			&share.SharedToType,
			&share.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shares, rows.Err()
}

// GetByUserID retrieves all shares by a user
func (r *PostgresShareRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Share, error) {
	query := `
		SELECT id, post_id, user_id, shared_to_type, created_at
		FROM post_shares
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*Share

	for rows.Next() {
		share := &Share{}
		err := rows.Scan(
			&share.ID,
			&share.PostID,
			&share.UserID,
			&share.SharedToType,
			&share.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shares, rows.Err()
}

// Delete removes a share from the database
func (r *PostgresShareRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM post_shares WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// GetShareCount returns the total number of shares for a post
func (r *PostgresShareRepository) GetShareCount(ctx context.Context, postID string) (int, error) {
	query := `SELECT COUNT(*) FROM post_shares WHERE post_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, postID).Scan(&count)
	return count, err
}

// HasUserShared checks if a user has already shared a post
func (r *PostgresShareRepository) HasUserShared(ctx context.Context, postID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post_shares WHERE post_id = $1 AND user_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, postID, userID).Scan(&exists)
	return exists, err
}


