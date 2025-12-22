package auth

import "context"

// OAuthRepository defines methods for OAuth account data access
type OAuthRepository interface {
	Create(ctx context.Context, account *OAuthAccount) error
	FindByProvider(ctx context.Context, provider, providerUserID string) (*OAuthAccount, error)
	FindByUserID(ctx context.Context, userID string) ([]*OAuthAccount, error)
	Delete(ctx context.Context, provider, providerUserID string) error
	DeleteByUserID(ctx context.Context, userID, provider string) error
}
