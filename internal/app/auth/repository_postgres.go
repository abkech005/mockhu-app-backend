package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository implements the UserRepository interface for PostgreSQL database.
// It handles all database operations related to User entities.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new instance of PostgresUserRepository.
// It takes a connection pool and returns a repository ready to interact with the database.
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// Create inserts a new user record into the database.
// It expects all required user fields to be populated before insertion.
// Returns an error if the operation fails (e.g., duplicate email, constraint violations).
// Empty strings for username, phone, first_name, last_name, and avatar_url are stored as NULL.
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (
		id, email, email_verified, phone, phone_verified, username, 
		first_name, last_name, dob, password_hash, avatar_url, 
		is_active, last_login_at, created_at, updated_at
	) VALUES (
		$1, $2, $3, NULLIF($4, ''), $5, NULLIF($6, ''), 
		NULLIF($7, ''), NULLIF($8, ''), $9, $10, NULLIF($11, ''), 
		$12, $13, $14, $15
	)`

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.EmailVerified, user.Phone, user.PhoneVerified,
		user.Username, user.FirstName, user.LastName, user.DOB, user.PasswordHash,
		user.AvatarURL, user.IsActive, user.LastLoginAt, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

// FindByID retrieves a user from the database by their unique ID.
// Returns the user if found, or an error if not found or query fails.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, 
		       COALESCE(first_name, '') as first_name, 
		       COALESCE(last_name, '') as last_name, 
		       dob, 
		       COALESCE(username, '') as username, 
		       password_hash,
		       email_verified, 
		       COALESCE(phone, '') as phone, 
		       phone_verified, 
		       COALESCE(avatar_url, '') as avatar_url, 
		       is_active,
		       created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DOB,
		&user.Username, &user.PasswordHash, &user.EmailVerified, &user.Phone,
		&user.PhoneVerified, &user.AvatarURL, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

// FindByEmail retrieves a user from the database by their email address.
// Returns the user if found, or an error if not found or query fails.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, 
		       COALESCE(first_name, '') as first_name, 
		       COALESCE(last_name, '') as last_name, 
		       dob, 
		       COALESCE(username, '') as username, 
		       password_hash,
		       email_verified, 
		       COALESCE(phone, '') as phone, 
		       phone_verified, 
		       COALESCE(avatar_url, '') as avatar_url, 
		       is_active,
		       created_at, updated_at, last_login_at
		FROM users WHERE email = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DOB,
		&user.Username, &user.PasswordHash, &user.EmailVerified, &user.Phone,
		&user.PhoneVerified, &user.AvatarURL, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// FindByPhone retrieves a user from the database by their phone number.
// Returns the user if found, or an error if not found or query fails.
func (r *PostgresUserRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	query := `
		SELECT id, email, 
		       COALESCE(first_name, '') as first_name, 
		       COALESCE(last_name, '') as last_name, 
		       dob, 
		       COALESCE(username, '') as username, 
		       password_hash,
		       email_verified, 
		       COALESCE(phone, '') as phone, 
		       phone_verified, 
		       COALESCE(avatar_url, '') as avatar_url, 
		       is_active,
		       created_at, updated_at, last_login_at
		FROM users WHERE phone = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DOB,
		&user.Username, &user.PasswordHash, &user.EmailVerified, &user.Phone,
		&user.PhoneVerified, &user.AvatarURL, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with phone %s not found", phone)
		}
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}

	return &user, nil
}

// FindByUsername retrieves a user from the database by their username.
// Returns the user if found, or an error if not found or query fails.
func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, email, 
		       COALESCE(first_name, '') as first_name, 
		       COALESCE(last_name, '') as last_name, 
		       dob, 
		       COALESCE(username, '') as username, 
		       password_hash,
		       email_verified, 
		       COALESCE(phone, '') as phone, 
		       phone_verified, 
		       COALESCE(avatar_url, '') as avatar_url, 
		       is_active,
		       created_at, updated_at, last_login_at
		FROM users WHERE username = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DOB,
		&user.Username, &user.PasswordHash, &user.EmailVerified, &user.Phone,
		&user.PhoneVerified, &user.AvatarURL, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &user, nil
}

// Update modifies an existing user record in the database.
// It updates all mutable fields based on the provided user object.
// Returns an error if the user doesn't exist or the operation fails.
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users SET
			email = $2, first_name = $3, last_name = $4, dob = $5,
			username = $6, password_hash = $7, email_verified = $8,
			phone = $9, phone_verified = $10, avatar_url = $11,
			is_active = $12, updated_at = $13, last_login_at = $14
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.FirstName, user.LastName, user.DOB,
		user.Username, user.PasswordHash, user.EmailVerified, user.Phone,
		user.PhoneVerified, user.AvatarURL, user.IsActive, user.UpdatedAt,
		user.LastLoginAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", user.ID)
	}

	return nil
}

// Delete removes a user record from the database by their ID.
// This is a hard delete operation. Consider soft delete for production systems.
// Returns an error if the user doesn't exist or the operation fails.
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}

// UpdateLastLogin updates the last_login_at timestamp for a user.
// This is typically called after successful authentication.
// Returns an error if the user doesn't exist or the operation fails.
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", userID)
	}

	return nil
}

// List retrieves a paginated list of users from the database.
// Parameters:
//   - limit: maximum number of users to return
//   - offset: number of users to skip (for pagination)
//
// Returns a slice of users ordered by creation date (newest first).
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	// TODO: Implement pagination query
	return nil, nil
}
