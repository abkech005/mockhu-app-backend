package education

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresEducationRepository implements EducationRepository using PostgreSQL
type PostgresEducationRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresEducationRepository creates a new PostgreSQL education repository
func NewPostgresEducationRepository(pool *pgxpool.Pool) *PostgresEducationRepository {
	return &PostgresEducationRepository{pool: pool}
}

// Create adds a new education entry to the database
func (r *PostgresEducationRepository) Create(ctx context.Context, education *Education) error {
	query := `
		INSERT INTO user_education (user_id, school, degree, field_of_study, location, start_year, end_year, current, logo_url, grade, activities, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		education.UserID,
		education.School,
		nullString(education.Degree),
		nullString(education.FieldOfStudy),
		nullString(education.Location),
		education.StartYear,
		education.EndYear,
		education.Current,
		nullString(education.LogoURL),
		nullString(education.Grade),
		nullString(education.Activities),
		nullString(education.Description),
	).Scan(&education.ID, &education.CreatedAt, &education.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create education entry: %w", err)
	}

	return nil
}

// FindByID retrieves an education entry by its ID
func (r *PostgresEducationRepository) FindByID(ctx context.Context, id string) (*Education, error) {
	query := `
		SELECT id, user_id, school, degree, field_of_study, location, start_year, end_year, current, logo_url, grade, activities, description, created_at, updated_at
		FROM user_education
		WHERE id = $1
	`

	var education Education
	var degree, fieldOfStudy, location, logoURL, grade, activities, description *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&education.ID,
		&education.UserID,
		&education.School,
		&degree,
		&fieldOfStudy,
		&location,
		&education.StartYear,
		&education.EndYear,
		&education.Current,
		&logoURL,
		&grade,
		&activities,
		&description,
		&education.CreatedAt,
		&education.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("education entry with id '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to find education entry: %w", err)
	}

	// Set nullable fields
	if degree != nil {
		education.Degree = *degree
	}
	if fieldOfStudy != nil {
		education.FieldOfStudy = *fieldOfStudy
	}
	if location != nil {
		education.Location = *location
	}
	if logoURL != nil {
		education.LogoURL = *logoURL
	}
	if grade != nil {
		education.Grade = *grade
	}
	if activities != nil {
		education.Activities = *activities
	}
	if description != nil {
		education.Description = *description
	}

	return &education, nil
}

// FindByUserID retrieves all education entries for a user
func (r *PostgresEducationRepository) FindByUserID(ctx context.Context, userID string) ([]Education, error) {
	query := `
		SELECT id, user_id, school, degree, field_of_study, location, start_year, end_year, current, logo_url, grade, activities, description, created_at, updated_at
		FROM user_education
		WHERE user_id = $1
		ORDER BY current DESC, start_year DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query education entries: %w", err)
	}
	defer rows.Close()

	var educations []Education
	for rows.Next() {
		var education Education
		var degree, fieldOfStudy, location, logoURL, grade, activities, description *string

		err := rows.Scan(
			&education.ID,
			&education.UserID,
			&education.School,
			&degree,
			&fieldOfStudy,
			&location,
			&education.StartYear,
			&education.EndYear,
			&education.Current,
			&logoURL,
			&grade,
			&activities,
			&description,
			&education.CreatedAt,
			&education.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan education entry: %w", err)
		}

		// Set nullable fields
		if degree != nil {
			education.Degree = *degree
		}
		if fieldOfStudy != nil {
			education.FieldOfStudy = *fieldOfStudy
		}
		if location != nil {
			education.Location = *location
		}
		if logoURL != nil {
			education.LogoURL = *logoURL
		}
		if grade != nil {
			education.Grade = *grade
		}
		if activities != nil {
			education.Activities = *activities
		}
		if description != nil {
			education.Description = *description
		}

		educations = append(educations, education)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating education entries: %w", err)
	}

	return educations, nil
}

// Update updates an education entry in the database
func (r *PostgresEducationRepository) Update(ctx context.Context, education *Education) error {
	query := `
		UPDATE user_education
		SET school = $1, degree = $2, field_of_study = $3, location = $4, start_year = $5, end_year = $6, 
		    current = $7, logo_url = $8, grade = $9, activities = $10, description = $11, updated_at = $12
		WHERE id = $13
		RETURNING updated_at
	`

	education.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query,
		education.School,
		nullString(education.Degree),
		nullString(education.FieldOfStudy),
		nullString(education.Location),
		education.StartYear,
		education.EndYear,
		education.Current,
		nullString(education.LogoURL),
		nullString(education.Grade),
		nullString(education.Activities),
		nullString(education.Description),
		education.UpdatedAt,
		education.ID,
	).Scan(&education.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("education entry with id '%s' not found", education.ID)
		}
		return fmt.Errorf("failed to update education entry: %w", err)
	}

	return nil
}

// Delete removes an education entry from the database
func (r *PostgresEducationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_education WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete education entry: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("education entry with id '%s' not found", id)
	}

	return nil
}

// FindCurrentEducation retrieves the current education for a user
func (r *PostgresEducationRepository) FindCurrentEducation(ctx context.Context, userID string) (*Education, error) {
	query := `
		SELECT id, user_id, school, degree, field_of_study, location, start_year, end_year, current, logo_url, grade, activities, description, created_at, updated_at
		FROM user_education
		WHERE user_id = $1 AND current = true
		LIMIT 1
	`

	var education Education
	var degree, fieldOfStudy, location, logoURL, grade, activities, description *string

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&education.ID,
		&education.UserID,
		&education.School,
		&degree,
		&fieldOfStudy,
		&location,
		&education.StartYear,
		&education.EndYear,
		&education.Current,
		&logoURL,
		&grade,
		&activities,
		&description,
		&education.CreatedAt,
		&education.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No current education
		}
		return nil, fmt.Errorf("failed to find current education: %w", err)
	}

	// Set nullable fields
	if degree != nil {
		education.Degree = *degree
	}
	if fieldOfStudy != nil {
		education.FieldOfStudy = *fieldOfStudy
	}
	if location != nil {
		education.Location = *location
	}
	if logoURL != nil {
		education.LogoURL = *logoURL
	}
	if grade != nil {
		education.Grade = *grade
	}
	if activities != nil {
		education.Activities = *activities
	}
	if description != nil {
		education.Description = *description
	}

	return &education, nil
}

// CountByUserID returns the count of education entries for a user
func (r *PostgresEducationRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM user_education WHERE user_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count education entries: %w", err)
	}

	return count, nil
}

// Helper function to convert empty strings to nil
func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
