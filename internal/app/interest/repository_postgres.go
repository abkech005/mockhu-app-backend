package interest

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresInterestRepository implements InterestRepository using PostgreSQL
type PostgresInterestRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresInterestRepository creates a new PostgreSQL interest repository
func NewPostgresInterestRepository(pool *pgxpool.Pool) *PostgresInterestRepository {
	return &PostgresInterestRepository{pool: pool}
}

// Create adds a new interest to the database
func (r *PostgresInterestRepository) Create(ctx context.Context, interest *Interest) error {
	query := `
		INSERT INTO interests (name, slug, category, icon)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	
	err := r.pool.QueryRow(ctx, query, interest.Name, interest.Slug, interest.Category, interest.Icon).
		Scan(&interest.ID, &interest.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create interest: %w", err)
	}
	
	return nil
}

// FindAll retrieves all interests from the database
func (r *PostgresInterestRepository) FindAll(ctx context.Context) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, created_at
		FROM interests
		ORDER BY category, name
	`
	
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query interests: %w", err)
	}
	defer rows.Close()
	
	var interests []Interest
	for rows.Next() {
		var interest Interest
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		interests = append(interests, interest)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating interests: %w", err)
	}
	
	return interests, nil
}

// FindBySlug retrieves an interest by its slug
func (r *PostgresInterestRepository) FindBySlug(ctx context.Context, slug string) (*Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, created_at
		FROM interests
		WHERE slug = $1
	`
	
	var interest Interest
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &interest.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("interest with slug '%s' not found", slug)
		}
		return nil, fmt.Errorf("failed to find interest: %w", err)
	}
	
	return &interest, nil
}

// FindBySlugs retrieves multiple interests by their slugs
func (r *PostgresInterestRepository) FindBySlugs(ctx context.Context, slugs []string) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, created_at
		FROM interests
		WHERE slug = ANY($1)
		ORDER BY name
	`
	
	rows, err := r.pool.Query(ctx, query, slugs)
	if err != nil {
		return nil, fmt.Errorf("failed to query interests by slugs: %w", err)
	}
	defer rows.Close()
	
	var interests []Interest
	for rows.Next() {
		var interest Interest
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		interests = append(interests, interest)
	}
	
	return interests, nil
}

// FindByCategory retrieves all interests in a specific category
func (r *PostgresInterestRepository) FindByCategory(ctx context.Context, category string) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, created_at
		FROM interests
		WHERE category = $1
		ORDER BY name
	`
	
	rows, err := r.pool.Query(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query interests by category: %w", err)
	}
	defer rows.Close()
	
	var interests []Interest
	for rows.Next() {
		var interest Interest
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		interests = append(interests, interest)
	}
	
	return interests, nil
}

// AddUserInterests adds multiple interests to a user
func (r *PostgresInterestRepository) AddUserInterests(ctx context.Context, userID string, interestIDs []string) error {
	// Use batch insert for better performance
	batch := &pgx.Batch{}
	
	for _, interestID := range interestIDs {
		query := `
			INSERT INTO user_interests (user_id, interest_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, interest_id) DO NOTHING
		`
		batch.Queue(query, userID, interestID)
	}
	
	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()
	
	for range interestIDs {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to add user interest: %w", err)
		}
	}
	
	return nil
}

// RemoveUserInterest removes an interest from a user
func (r *PostgresInterestRepository) RemoveUserInterest(ctx context.Context, userID string, interestID string) error {
	query := `
		DELETE FROM user_interests
		WHERE user_id = $1 AND interest_id = $2
	`
	
	result, err := r.pool.Exec(ctx, query, userID, interestID)
	if err != nil {
		return fmt.Errorf("failed to remove user interest: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user interest not found")
	}
	
	return nil
}

// GetUserInterests retrieves all interests for a user
func (r *PostgresInterestRepository) GetUserInterests(ctx context.Context, userID string) ([]Interest, error) {
	query := `
		SELECT i.id, i.name, i.slug, i.category, i.icon, i.created_at
		FROM interests i
		INNER JOIN user_interests ui ON i.id = ui.interest_id
		WHERE ui.user_id = $1
		ORDER BY i.category, i.name
	`
	
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user interests: %w", err)
	}
	defer rows.Close()
	
	var interests []Interest
	for rows.Next() {
		var interest Interest
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		interests = append(interests, interest)
	}
	
	return interests, nil
}

// ReplaceUserInterests removes all user interests and adds new ones
func (r *PostgresInterestRepository) ReplaceUserInterests(ctx context.Context, userID string, interestIDs []string) error {
	// Use transaction to ensure atomicity
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Delete all existing interests
	deleteQuery := `DELETE FROM user_interests WHERE user_id = $1`
	_, err = tx.Exec(ctx, deleteQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to delete existing interests: %w", err)
	}
	
	// Insert new interests
	if len(interestIDs) > 0 {
		for _, interestID := range interestIDs {
			insertQuery := `
				INSERT INTO user_interests (user_id, interest_id)
				VALUES ($1, $2)
			`
			_, err = tx.Exec(ctx, insertQuery, userID, interestID)
			if err != nil {
				return fmt.Errorf("failed to insert new interest: %w", err)
			}
		}
	}
	
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// UserHasInterest checks if a user has a specific interest
func (r *PostgresInterestRepository) UserHasInterest(ctx context.Context, userID string, interestID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_interests
			WHERE user_id = $1 AND interest_id = $2
		)
	`
	
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, interestID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user interest: %w", err)
	}
	
	return exists, nil
}

// CountByCategory returns the count of interests per category
func (r *PostgresInterestRepository) CountByCategory(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT category, COUNT(*) as count
		FROM interests
		GROUP BY category
		ORDER BY category
	`
	
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count interests by category: %w", err)
	}
	defer rows.Close()
	
	counts := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category count: %w", err)
		}
		counts[category] = count
	}
	
	return counts, nil
}

// CountUserInterests returns the count of interests for a user
func (r *PostgresInterestRepository) CountUserInterests(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_interests
		WHERE user_id = $1
	`
	
	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user interests: %w", err)
	}
	
	return count, nil
}

