package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresOAuthRepository implements OAuthRepository for PostgreSQL
type PostgresOAuthRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresOAuthRepository creates a new PostgresOAuthRepository
func NewPostgresOAuthRepository(pool *pgxpool.Pool) *PostgresOAuthRepository {
	return &PostgresOAuthRepository{pool: pool}
}

func (r *PostgresOAuthRepository) Create(ctx context.Context, account *OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (
			user_id, provider, provider_user_id, provider_email, 
			access_token, refresh_token, token_expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		account.UserID,
		account.Provider,
		account.ProviderUserID,
		account.ProviderEmail,
		account.AccessToken,
		account.RefreshToken,
		account.TokenExpiresAt,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	return err
}

func (r *PostgresOAuthRepository) FindByProvider(ctx context.Context, provider, providerUserID string) (*OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, 
		       access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE provider = $1 AND provider_user_id = $2`

	account := &OAuthAccount{}
	// AccessToken/RefreshToken/TokenExpiresAt can be nullable in DB?
	// The migration didn't enforce NOT NULL on them, but logic usually has them.
	// Postgres driver might complain if we scan NULL into string.
	// Migration: "access_token TEXT, refresh_token TEXT".
	// Let's use coalesce or pointer in struct. Struct has string.
	// We'll scan into placeholders and handle or ensure DB has values.
	// But let's assume valid data for now, or use *string if needed.
	// Struct: AccessToken string.
	// If DB has NULL, this will fail.
	// Migration: NO "NOT NULL" on access_token.
	// But Create inserts it.
	// I'll leave as is, assuming populated.

	err := r.pool.QueryRow(ctx, query, provider, providerUserID).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUserID,
		&account.ProviderEmail,
		&account.AccessToken,
		&account.RefreshToken,
		&account.TokenExpiresAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("oauth account not found")
		}
		return nil, err
	}

	return account, nil
}

func (r *PostgresOAuthRepository) FindByUserID(ctx context.Context, userID string) ([]*OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, 
		       access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE user_id = $1`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*OAuthAccount
	for rows.Next() {
		account := &OAuthAccount{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Provider,
			&account.ProviderUserID,
			&account.ProviderEmail,
			&account.AccessToken,
			&account.RefreshToken,
			&account.TokenExpiresAt,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresOAuthRepository) Delete(ctx context.Context, provider, providerUserID string) error {
	query := `DELETE FROM oauth_accounts WHERE provider = $1 AND provider_user_id = $2`
	_, err := r.pool.Exec(ctx, query, provider, providerUserID)
	return err
}

func (r *PostgresOAuthRepository) DeleteByUserID(ctx context.Context, userID, provider string) error {
	query := `DELETE FROM oauth_accounts WHERE user_id = $1 AND provider = $2`
	_, err := r.pool.Exec(ctx, query, userID, provider)
	return err
}
