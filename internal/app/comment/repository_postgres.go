package comment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresCommentRepository implements CommentRepository for PostgreSQL
type PostgresCommentRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresCommentRepository creates a new PostgreSQL comment repository
func NewPostgresCommentRepository(pool *pgxpool.Pool) CommentRepository {
	return &PostgresCommentRepository{pool: pool}
}

// Create inserts a new comment into the database
func (r *PostgresCommentRepository) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO post_comments (post_id, user_id, parent_comment_id, content, is_anonymous)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Content,
		comment.IsAnonymous,
	).Scan(&comment.ID, &comment.IsActive, &comment.CreatedAt, &comment.UpdatedAt)

	return err
}

// GetByID retrieves a single comment by its ID
func (r *PostgresCommentRepository) GetByID(ctx context.Context, id string) (*Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_comment_id, content, is_anonymous, 
		       is_active, created_at, updated_at
		FROM post_comments
		WHERE id = $1 AND is_active = true
	`

	comment := &Comment{}
	var parentCommentID *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&parentCommentID,
		&comment.Content,
		&comment.IsAnonymous,
		&comment.IsActive,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	comment.ParentCommentID = parentCommentID
	return comment, nil
}

// GetByPostID retrieves all top-level comments for a post (no parent)
func (r *PostgresCommentRepository) GetByPostID(ctx context.Context, postID string, limit, offset int) ([]*Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_comment_id, content, is_anonymous,
		       is_active, created_at, updated_at
		FROM post_comments
		WHERE post_id = $1 AND parent_comment_id IS NULL AND is_active = true
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		comment := &Comment{}
		var parentCommentID *string

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&parentCommentID,
			&comment.Content,
			&comment.IsAnonymous,
			&comment.IsActive,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		comment.ParentCommentID = parentCommentID
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// GetReplies retrieves all replies to a parent comment
func (r *PostgresCommentRepository) GetReplies(ctx context.Context, parentCommentID string, limit, offset int) ([]*Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_comment_id, content, is_anonymous,
		       is_active, created_at, updated_at
		FROM post_comments
		WHERE parent_comment_id = $1 AND is_active = true
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, parentCommentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		comment := &Comment{}
		var parentCommentID *string

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&parentCommentID,
			&comment.Content,
			&comment.IsAnonymous,
			&comment.IsActive,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		comment.ParentCommentID = parentCommentID
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// Update modifies an existing comment
func (r *PostgresCommentRepository) Update(ctx context.Context, comment *Comment) error {
	query := `
		UPDATE post_comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND is_active = true
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query, comment.Content, comment.ID, comment.UserID).
		Scan(&comment.UpdatedAt)

	return err
}

// Delete performs a soft delete on a comment
func (r *PostgresCommentRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE post_comments
		SET is_active = false, updated_at = NOW()
		WHERE id = $1 AND is_active = true
	`

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

// GetCommentCount returns the total number of comments for a post
func (r *PostgresCommentRepository) GetCommentCount(ctx context.Context, postID string) (int, error) {
	query := `SELECT COUNT(*) FROM post_comments WHERE post_id = $1 AND is_active = true`

	var count int
	err := r.pool.QueryRow(ctx, query, postID).Scan(&count)
	return count, err
}

// GetReplyCount returns the total number of replies to a comment
func (r *PostgresCommentRepository) GetReplyCount(ctx context.Context, parentCommentID string) (int, error) {
	query := `SELECT COUNT(*) FROM post_comments WHERE parent_comment_id = $1 AND is_active = true`

	var count int
	err := r.pool.QueryRow(ctx, query, parentCommentID).Scan(&count)
	return count, err
}

