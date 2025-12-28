package location

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresLocationRepository implements LocationRepository using PostgreSQL
type PostgresLocationRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresLocationRepository creates a new PostgreSQL location repository
func NewPostgresLocationRepository(pool *pgxpool.Pool) *PostgresLocationRepository {
	return &PostgresLocationRepository{pool: pool}
}

// Create adds a new location to the database
func (r *PostgresLocationRepository) Create(ctx context.Context, location *Location) error {
	query := `
		INSERT INTO locations (city, country, used_by_count)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, location.City, location.Country, location.UsedByCount).
		Scan(&location.ID, &location.CreatedAt, &location.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

// FindByID retrieves a location by its ID
func (r *PostgresLocationRepository) FindByID(ctx context.Context, id string) (*Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		WHERE id = $1
	`

	var location Location
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("location with id '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to find location: %w", err)
	}

	return &location, nil
}

// FindAll retrieves all locations
func (r *PostgresLocationRepository) FindAll(ctx context.Context) ([]Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		ORDER BY used_by_count DESC, city
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query locations: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

// Update updates a location
func (r *PostgresLocationRepository) Update(ctx context.Context, location *Location) error {
	query := `
		UPDATE locations
		SET city = $1, country = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at
	`

	location.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query, location.City, location.Country, location.UpdatedAt, location.ID).
		Scan(&location.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("location with id '%s' not found", location.ID)
		}
		return fmt.Errorf("failed to update location: %w", err)
	}

	return nil
}

// Delete removes a location
func (r *PostgresLocationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM locations WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("location with id '%s' not found", id)
	}

	return nil
}

// Search searches locations by city or country (autocomplete)
func (r *PostgresLocationRepository) Search(ctx context.Context, query string) ([]Location, error) {
	sqlQuery := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		WHERE city ILIKE $1 OR country ILIKE $1
		ORDER BY used_by_count DESC, city
		LIMIT 20
	`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, sqlQuery, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search locations: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

// FindByCity finds all locations in a city
func (r *PostgresLocationRepository) FindByCity(ctx context.Context, city string) ([]Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		WHERE city ILIKE $1
		ORDER BY used_by_count DESC
	`

	rows, err := r.pool.Query(ctx, query, city)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations by city: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

// FindByCountry finds all locations in a country
func (r *PostgresLocationRepository) FindByCountry(ctx context.Context, country string) ([]Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		WHERE country ILIKE $1
		ORDER BY used_by_count DESC, city
	`

	rows, err := r.pool.Query(ctx, query, country)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations by country: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

// FindByCityAndCountry finds a location by exact city and country
func (r *PostgresLocationRepository) FindByCityAndCountry(ctx context.Context, city, country string) (*Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		WHERE city ILIKE $1 AND country ILIKE $2
	`

	var location Location
	err := r.pool.QueryRow(ctx, query, city, country).Scan(
		&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find location: %w", err)
	}

	return &location, nil
}

// IncrementUsedByCount increments the used_by_count
func (r *PostgresLocationRepository) IncrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE locations
		SET used_by_count = used_by_count + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("location with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to increment used_by_count: %w", err)
	}

	return count, nil
}

// DecrementUsedByCount decrements the used_by_count (min 0)
func (r *PostgresLocationRepository) DecrementUsedByCount(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE locations
		SET used_by_count = GREATEST(used_by_count - 1, 0), updated_at = NOW()
		WHERE id = $1
		RETURNING used_by_count
	`

	var count int
	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("location with id '%s' not found", id)
		}
		return 0, fmt.Errorf("failed to decrement used_by_count: %w", err)
	}

	return count, nil
}

// GetMostUsedLocations returns the most used locations
func (r *PostgresLocationRepository) GetMostUsedLocations(ctx context.Context, limit int) ([]Location, error) {
	query := `
		SELECT id, city, country, used_by_count, created_at, updated_at
		FROM locations
		ORDER BY used_by_count DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query most used locations: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.City, &location.Country, &location.UsedByCount, &location.CreatedAt, &location.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return locations, nil
}

// CountByCountry returns count of locations by country
func (r *PostgresLocationRepository) CountByCountry(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT country, COUNT(*) as count
		FROM locations
		GROUP BY country
		ORDER BY count DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count by country: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var country string
		var count int
		if err := rows.Scan(&country, &count); err != nil {
			return nil, fmt.Errorf("failed to scan country count: %w", err)
		}
		counts[country] = count
	}

	return counts, nil
}
