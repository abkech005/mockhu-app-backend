package title

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresTitleRepository implements TitleRepository using PostgreSQL
type PostgresTitleRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresTitleRepository creates a new PostgreSQL title repository
func NewPostgresTitleRepository(pool *pgxpool.Pool) *PostgresTitleRepository {
	return &PostgresTitleRepository{pool: pool}
}

// Create adds a new title to the database
func (r *PostgresTitleRepository) Create(ctx context.Context, title *Title) error {
	query := `
		INSERT INTO titles (name, description, defined_by, used_by_count)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, title.Name, title.Description, title.DefinedBy, title.UsedByCount).
		Scan(&title.ID, &title.CreatedAt, &title.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create title: %w", err)
	}

	return nil
}

// FindAll retrieves all titles from the database
func (r *PostgresTitleRepository) FindAll(ctx context.Context) ([]Title, error) {
	query := `
		SELECT id, name, description, defined_by, used_by_count, created_at, updated_at
		FROM titles
		ORDER BY defined_by, name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query titles: %w", err)
	}
	defer rows.Close()

	var titles []Title
	for rows.Next() {
		var title Title
		var description *string
		err := rows.Scan(&title.ID, &title.Name, &description, &title.DefinedBy, &title.UsedByCount, &title.CreatedAt, &title.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan title: %w", err)
		}
		if description != nil {
			title.Description = *description
		}
		titles = append(titles, title)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating titles: %w", err)
	}

	return titles, nil
}

// FindByID retrieves a title by its ID
func (r *PostgresTitleRepository) FindByID(ctx context.Context, id string) (*Title, error) {
	query := `
		SELECT id, name, description, defined_by, used_by_count, created_at, updated_at
		FROM titles
		WHERE id = $1
	`

	var title Title
	var description *string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&title.ID, &title.Name, &description, &title.DefinedBy, &title.UsedByCount, &title.CreatedAt, &title.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("title with id '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to find title: %w", err)
	}

	if description != nil {
		title.Description = *description
	}

	return &title, nil
}

// FindByName retrieves a title by its name
func (r *PostgresTitleRepository) FindByName(ctx context.Context, name string) (*Title, error) {
	query := `
		SELECT id, name, description, defined_by, used_by_count, created_at, updated_at
		FROM titles
		WHERE name = $1
	`

	var title Title
	var description *string
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&title.ID, &title.Name, &description, &title.DefinedBy, &title.UsedByCount, &title.CreatedAt, &title.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("title with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to find title: %w", err)
	}

	if description != nil {
		title.Description = *description
	}

	return &title, nil
}

// FindByDefinedBy retrieves all titles defined by admin or user
func (r *PostgresTitleRepository) FindByDefinedBy(ctx context.Context, definedBy string) ([]Title, error) {
	query := `
		SELECT id, name, description, defined_by, used_by_count, created_at, updated_at
		FROM titles
		WHERE defined_by = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, definedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to query titles by defined_by: %w", err)
	}
	defer rows.Close()

	var titles []Title
	for rows.Next() {
		var title Title
		var description *string
		err := rows.Scan(&title.ID, &title.Name, &description, &title.DefinedBy, &title.UsedByCount, &title.CreatedAt, &title.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan title: %w", err)
		}
		if description != nil {
			title.Description = *description
		}
		titles = append(titles, title)
	}

	return titles, nil
}

// Update updates a title in the database
func (r *PostgresTitleRepository) Update(ctx context.Context, title *Title) error {
	query := `
		UPDATE titles
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at
	`

	title.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query, title.Name, title.Description, title.UpdatedAt, title.ID).
		Scan(&title.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("title with id '%s' not found", title.ID)
		}
		return fmt.Errorf("failed to update title: %w", err)
	}

	return nil
}

// Delete removes a title from the database
func (r *PostgresTitleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM titles WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete title: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("title with id '%s' not found", id)
	}

	return nil
}

// IncrementUsedByCount increments the used_by_count for a title
func (r *PostgresTitleRepository) IncrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE titles
		SET used_by_count = used_by_count + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("title with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to increment used_by_count: %w", err)
	}

	return count, nil
}

// DecrementUsedByCount decrements the used_by_count for a title (min 0)
func (r *PostgresTitleRepository) DecrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE titles
		SET used_by_count = GREATEST(used_by_count - 1, 0), updated_at = NOW()
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("title with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to decrement used_by_count: %w", err)
	}

	return count, nil
}

// CountByDefinedBy returns the count of titles per defined_by
func (r *PostgresTitleRepository) CountByDefinedBy(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT defined_by, COUNT(*) as count
		FROM titles
		GROUP BY defined_by
		ORDER BY defined_by
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count titles by defined_by: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var definedBy string
		var count int
		err := rows.Scan(&definedBy, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[definedBy] = count
	}

	return counts, nil
}

// GetMostUsedTitles returns the most used titles
func (r *PostgresTitleRepository) GetMostUsedTitles(ctx context.Context, limit int) ([]Title, error) {
	query := `
		SELECT id, name, description, defined_by, used_by_count, created_at, updated_at
		FROM titles
		ORDER BY used_by_count DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query most used titles: %w", err)
	}
	defer rows.Close()

	var titles []Title
	for rows.Next() {
		var title Title
		var description *string
		err := rows.Scan(&title.ID, &title.Name, &description, &title.DefinedBy, &title.UsedByCount, &title.CreatedAt, &title.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan title: %w", err)
		}
		if description != nil {
			title.Description = *description
		}
		titles = append(titles, title)
	}

	return titles, nil
}
