package postfeed

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements Repository for PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new postfeed
func (r *PostgresRepository) Create(ctx context.Context, p *Postfeed) error {
	query := `
		INSERT INTO postfeeds (user_id, type, title, content, tags, is_anonymous, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		p.UserID,
		p.Type,
		p.Title,
		p.Content,
		p.Tags,
		p.IsAnonymous,
		p.Metadata,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// GetByID retrieves a postfeed by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Postfeed, error) {
	query := `
		SELECT id, user_id, type, title, content, tags, is_anonymous, is_active,
		       metadata, view_count, like_count, comment_count, share_count,
		       created_at, updated_at
		FROM postfeeds
		WHERE id = $1 AND is_active = true
	`

	p := &Postfeed{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.Type, &p.Title, &p.Content, &p.Tags,
		&p.IsAnonymous, &p.IsActive, &p.Metadata, &p.ViewCount,
		&p.LikeCount, &p.CommentCount, &p.ShareCount,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("postfeed not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get postfeed: %w", err)
	}
	return p, nil
}

// Update updates a postfeed
func (r *PostgresRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	var setParts []string
	var args []interface{}
	argPos := 1

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argPos))
		args = append(args, value)
		argPos++
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	query := fmt.Sprintf("UPDATE postfeeds SET %s WHERE id = $%d AND is_active = true",
		strings.Join(setParts, ", "), argPos)
	args = append(args, id)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update postfeed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("postfeed not found")
	}

	return nil
}

// Delete soft deletes a postfeed
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE postfeeds SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete postfeed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("postfeed not found")
	}
	return nil
}

// List retrieves postfeeds with filters
func (r *PostgresRepository) List(ctx context.Context, filter ListFilter) ([]Postfeed, int, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	conditions = append(conditions, "is_active = true")

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, filter.Type)
		argPos++
	}

	if filter.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, filter.UserID)
		argPos++
	}

	if filter.Tag != "" {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(tags)", argPos))
		args = append(args, filter.Tag)
		argPos++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM postfeeds WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count postfeeds: %w", err)
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, user_id, type, title, content, tags, is_anonymous, is_active,
		       metadata, view_count, like_count, comment_count, share_count,
		       created_at, updated_at
		FROM postfeeds
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list postfeeds: %w", err)
	}
	defer rows.Close()

	var postfeeds []Postfeed
	for rows.Next() {
		var p Postfeed
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Type, &p.Title, &p.Content, &p.Tags,
			&p.IsAnonymous, &p.IsActive, &p.Metadata, &p.ViewCount,
			&p.LikeCount, &p.CommentCount, &p.ShareCount,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan postfeed: %w", err)
		}
		postfeeds = append(postfeeds, p)
	}

	return postfeeds, total, nil
}

// GetByUserID retrieves postfeeds by user
func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]Postfeed, int, error) {
	return r.List(ctx, ListFilter{UserID: userID, Limit: limit, Offset: offset})
}

// GetByType retrieves postfeeds by type
func (r *PostgresRepository) GetByType(ctx context.Context, postType string, limit, offset int) ([]Postfeed, int, error) {
	return r.List(ctx, ListFilter{Type: postType, Limit: limit, Offset: offset})
}

// IncrementViewCount increments view count
func (r *PostgresRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `UPDATE postfeeds SET view_count = view_count + 1 WHERE id = $1 AND is_active = true`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdateEngagementCounts updates engagement counts
func (r *PostgresRepository) UpdateEngagementCounts(ctx context.Context, id string, likes, comments, shares int) error {
	query := `
		UPDATE postfeeds 
		SET like_count = like_count + $2,
		    comment_count = comment_count + $3,
		    share_count = share_count + $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
	`
	_, err := r.db.Exec(ctx, query, id, likes, comments, shares)
	return err
}

// Helper to unmarshal metadata
func UnmarshalMetadata[T any](data json.RawMessage) (*T, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
