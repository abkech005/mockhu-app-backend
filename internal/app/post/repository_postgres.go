package post

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPostRepository implements PostRepository for PostgreSQL
type PostgresPostRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPostRepository(pool *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{pool: pool}
}

func (r *PostgresPostRepository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, content, images, is_anonymous)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_active, view_count, created_at, updated_at
	`

	if post.Images == nil {
		post.Images = []string{}
	}

	err := r.pool.QueryRow(ctx, query, post.UserID, post.Content, post.Images, post.IsAnonymous).
		Scan(&post.ID, &post.IsActive, &post.ViewCount, &post.CreatedAt, &post.UpdatedAt)

	return err
}

// GetByID retrieves a single post by its ID
func (r *PostgresPostRepository) GetByID(ctx context.Context, id string) (*Post, error) {
	query := `
		SELECT id, user_id, content, images, is_anonymous, is_active, 
		       view_count, created_at, updated_at
		FROM posts
		WHERE id = $1 AND is_active = true
	`

	post := &Post{}
	var images []string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&images,
		&post.IsAnonymous,
		&post.IsActive,
		&post.ViewCount,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	post.Images = images
	return post, nil
}

// GetByUserID retrieves all active posts by a specific user
func (r *PostgresPostRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Post, error) {
	query := `
		SELECT id, user_id, content, images, is_anonymous, is_active,
		       view_count, created_at, updated_at
		FROM posts
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		post := &Post{}
		var images []string

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&images,
			&post.IsAnonymous,
			&post.IsActive,
			&post.ViewCount,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		post.Images = images
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// Update modifies an existing post
func (r *PostgresPostRepository) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET content = $1, images = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4 AND is_active = true
		RETURNING updated_at
	`

	if post.Images == nil {
		post.Images = []string{}
	}

	err := r.pool.QueryRow(ctx, query, post.Content, post.Images, post.ID, post.UserID).
		Scan(&post.UpdatedAt)

	return err
}

// Delete performs a soft delete on a post
func (r *PostgresPostRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE posts
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

// GetFeed retrieves posts from users that the current user follows
func (r *PostgresPostRepository) GetFeed(ctx context.Context, userID string, limit, offset int) ([]*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.images, p.is_anonymous, p.is_active,
		       p.view_count, p.created_at, p.updated_at
		FROM posts p
		WHERE p.user_id IN (
			SELECT following_id FROM user_follows WHERE follower_id = $1
		)
		AND p.is_active = true
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		post := &Post{}
		var images []string

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&images,
			&post.IsAnonymous,
			&post.IsActive,
			&post.ViewCount,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		post.Images = images
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// AddReaction adds a fire reaction to a post
func (r *PostgresPostRepository) AddReaction(ctx context.Context, reaction *Reaction) error {
	query := `
		INSERT INTO post_reactions (post_id, user_id, reaction_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, user_id) DO NOTHING
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, reaction.PostID, reaction.UserID, reaction.ReactionType).
		Scan(&reaction.ID, &reaction.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // ON CONFLICT triggered, reaction already exists
		}
		return err
	}

	return nil
}

// RemoveReaction removes a user's reaction from a post
func (r *PostgresPostRepository) RemoveReaction(ctx context.Context, postID, userID string) error {
	query := `
		DELETE FROM post_reactions
		WHERE post_id = $1 AND user_id = $2
	`

	result, err := r.pool.Exec(ctx, query, postID, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// GetReactions retrieves all reactions for a post
func (r *PostgresPostRepository) GetReactions(ctx context.Context, postID string, limit, offset int) ([]*Reaction, error) {
	query := `
		SELECT id, post_id, user_id, reaction_type, created_at
		FROM post_reactions
		WHERE post_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*Reaction

	for rows.Next() {
		reaction := &Reaction{}
		err := rows.Scan(
			&reaction.ID,
			&reaction.PostID,
			&reaction.UserID,
			&reaction.ReactionType,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}

	return reactions, rows.Err()
}

// GetReactionCount returns the total number of reactions for a post
func (r *PostgresPostRepository) GetReactionCount(ctx context.Context, postID string) (int, error) {
	query := `SELECT COUNT(*) FROM post_reactions WHERE post_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, postID).Scan(&count)
	return count, err
}

// HasUserReacted checks if a specific user has reacted to a post
func (r *PostgresPostRepository) HasUserReacted(ctx context.Context, postID, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM post_reactions
			WHERE post_id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, postID, userID).Scan(&exists)
	return exists, err
}
