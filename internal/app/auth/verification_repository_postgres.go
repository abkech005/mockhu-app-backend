package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresVerificationRepository implements the VerificationRepository interface for PostgreSQL database.
// It handles all database operations related to VerificationCode entities.
type PostgresVerificationRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresVerificationRepository creates a new instance of PostgresVerificationRepository.
// It takes a connection pool and returns a repository ready to interact with the database.
func NewPostgresVerificationRepository(pool *pgxpool.Pool) *PostgresVerificationRepository {
	return &PostgresVerificationRepository{pool: pool}
}

// Create inserts a new verification code into the database.
// It expects all required fields to be populated before insertion.
func (r *PostgresVerificationRepository) Create(ctx context.Context, verification *VerificationCode) error {
	query := `
		INSERT INTO verification_codes (
			id, user_id, code, type, contact, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		verification.ID,
		verification.UserID,
		verification.Code,
		verification.Type,
		verification.Contact,
		verification.ExpiresAt,
		verification.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create verification code: %w", err)
	}
	return nil
}

// FindByCodeAndType finds a verification code by code and type.
// Returns ErrNoRows if no matching code is found.
func (r *PostgresVerificationRepository) FindByCodeAndType(ctx context.Context, code string, verificationType string) (*VerificationCode, error) {
	query := `
		SELECT id, user_id, code, type, contact, is_active, used_at, expires_at, created_at 
		FROM verification_codes 
		WHERE code = $1 AND type = $2 AND is_active = true`

	var verification VerificationCode
	err := r.pool.QueryRow(ctx, query, code, verificationType).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Code,
		&verification.Type,
		&verification.Contact,
		&verification.IsActive,
		&verification.UsedAt,
		&verification.ExpiresAt,
		&verification.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("verification code not found")
		}
		return nil, fmt.Errorf("failed to find verification code: %w", err)
	}

	return &verification, nil
}

// FindActiveByContactAndType finds the latest active (unused and non-expired) verification code for a contact.
// Returns ErrNoRows if no active code is found.
func (r *PostgresVerificationRepository) FindActiveByContactAndType(ctx context.Context, contact string, verificationType string) (*VerificationCode, error) {
	query := `
		SELECT id, user_id, code, type, contact, used_at, expires_at, created_at 
		FROM verification_codes 
		WHERE contact = $1 AND type = $2 AND used_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1`

	var verification VerificationCode
	err := r.pool.QueryRow(ctx, query, contact, verificationType).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Code,
		&verification.Type,
		&verification.Contact,
		&verification.UsedAt,
		&verification.ExpiresAt,
		&verification.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no active verification code found")
		}
		return nil, fmt.Errorf("failed to find active verification code: %w", err)
	}

	return &verification, nil
}

// MarkAsUsed marks a verification code as used by setting the used_at timestamp and deactivating it.
// This prevents the code from being reused.
func (r *PostgresVerificationRepository) MarkAsUsed(ctx context.Context, id string) error {
	query := `UPDATE verification_codes SET used_at = $1, is_active = false WHERE id = $2`

	result, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark verification code as used: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("verification code not found")
	}

	return nil
}

// DeactivatePreviousCodes deactivates all previous active verification codes for a user and type.
// This is typically called when generating a new verification code to invalidate old ones.
// More efficient than marking as used - just sets is_active to false.
func (r *PostgresVerificationRepository) DeactivatePreviousCodes(ctx context.Context, userID string, verificationType string) error {
	query := `
		UPDATE verification_codes 
		SET is_active = false 
		WHERE user_id = $1 AND type = $2 AND is_active = true`

	_, err := r.pool.Exec(ctx, query, userID, verificationType)
	if err != nil {
		return fmt.Errorf("failed to deactivate previous codes: %w", err)
	}

	return nil
}

// CleanupExpired deletes all expired or inactive verification codes from the database.
// Returns the number of deleted rows.
// This should be called periodically by a background job to keep the database clean.
func (r *PostgresVerificationRepository) CleanupExpired(ctx context.Context) (int64, error) {
	query := `
		DELETE FROM verification_codes 
		WHERE expires_at < NOW() 
		   OR (is_active = false AND created_at < NOW() - INTERVAL '7 days')`

	result, err := r.pool.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired codes: %w", err)
	}

	return result.RowsAffected(), nil
}
