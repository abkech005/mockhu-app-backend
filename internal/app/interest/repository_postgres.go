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
		INSERT INTO interests (name, slug, category, icon, defined_by, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, used_by_count, created_at
	`

	err := r.pool.QueryRow(ctx, query, interest.Name, interest.Slug, interest.Category, interest.Icon, interest.DefinedBy, nullString(interest.Description)).
		Scan(&interest.ID, &interest.UsedByCount, &interest.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create interest: %w", err)
	}

	return nil
}

// Helper to convert empty string to nil
func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// FindAll retrieves all interests from the database
func (r *PostgresInterestRepository) FindAll(ctx context.Context) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
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
		var description, definedBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &definedBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if definedBy != nil {
			interest.DefinedBy = *definedBy
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
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
		FROM interests
		WHERE slug = $1
	`

	var interest Interest
	var description, definedBy *string
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &definedBy, &interest.UsedByCount, &description, &interest.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("interest with slug '%s' not found", slug)
		}
		return nil, fmt.Errorf("failed to find interest: %w", err)
	}

	if description != nil {
		interest.Description = *description
	}
	if definedBy != nil {
		interest.DefinedBy = *definedBy
	}

	return &interest, nil
}

// FindBySlugs retrieves multiple interests by their slugs
func (r *PostgresInterestRepository) FindBySlugs(ctx context.Context, slugs []string) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
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
		var description, definedBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &definedBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if definedBy != nil {
			interest.DefinedBy = *definedBy
		}
		interests = append(interests, interest)
	}

	return interests, nil
}

// FindByCategory retrieves all interests in a specific category
func (r *PostgresInterestRepository) FindByCategory(ctx context.Context, category string) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
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
		var description, definedBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &definedBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if definedBy != nil {
			interest.DefinedBy = *definedBy
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
		SELECT i.id, i.name, i.slug, i.category, i.icon, i.defined_by, i.used_by_count, i.description, i.created_at
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
		var description, definedBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &definedBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if definedBy != nil {
			interest.DefinedBy = *definedBy
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

// FindByDefinedBy retrieves all interests defined by admin or user
func (r *PostgresInterestRepository) FindByDefinedBy(ctx context.Context, definedBy string) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
		FROM interests
		WHERE defined_by = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, definedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to query interests by defined_by: %w", err)
	}
	defer rows.Close()

	var interests []Interest
	for rows.Next() {
		var interest Interest
		var description, defBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &defBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if defBy != nil {
			interest.DefinedBy = *defBy
		}
		interests = append(interests, interest)
	}

	return interests, nil
}

// IncrementUsedByCount increments the used_by_count for an interest
func (r *PostgresInterestRepository) IncrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE interests
		SET used_by_count = used_by_count + 1
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("interest with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to increment used_by_count: %w", err)
	}

	return count, nil
}

// DecrementUsedByCount decrements the used_by_count for an interest (min 0)
func (r *PostgresInterestRepository) DecrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE interests
		SET used_by_count = GREATEST(used_by_count - 1, 0)
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("interest with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to decrement used_by_count: %w", err)
	}

	return count, nil
}

// GetMostUsedInterests returns the most used interests
func (r *PostgresInterestRepository) GetMostUsedInterests(ctx context.Context, limit int) ([]Interest, error) {
	query := `
		SELECT id, name, slug, category, icon, defined_by, used_by_count, description, created_at
		FROM interests
		ORDER BY used_by_count DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query most used interests: %w", err)
	}
	defer rows.Close()

	var interests []Interest
	for rows.Next() {
		var interest Interest
		var description, defBy *string
		err := rows.Scan(&interest.ID, &interest.Name, &interest.Slug, &interest.Category, &interest.Icon, &defBy, &interest.UsedByCount, &description, &interest.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		if description != nil {
			interest.Description = *description
		}
		if defBy != nil {
			interest.DefinedBy = *defBy
		}
		interests = append(interests, interest)
	}

	return interests, nil
}
