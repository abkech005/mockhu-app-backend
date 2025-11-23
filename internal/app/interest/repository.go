package interest

import "context"

// InterestRepository defines methods for interest data access
type InterestRepository interface {
	// Interest CRUD
	Create(ctx context.Context, interest *Interest) error
	FindAll(ctx context.Context) ([]Interest, error)
	FindBySlug(ctx context.Context, slug string) (*Interest, error)
	FindBySlugs(ctx context.Context, slugs []string) ([]Interest, error)
	FindByCategory(ctx context.Context, category string) ([]Interest, error)
	
	// User Interest Management
	AddUserInterests(ctx context.Context, userID string, interestIDs []string) error
	RemoveUserInterest(ctx context.Context, userID string, interestID string) error
	GetUserInterests(ctx context.Context, userID string) ([]Interest, error)
	ReplaceUserInterests(ctx context.Context, userID string, interestIDs []string) error
	UserHasInterest(ctx context.Context, userID string, interestID string) (bool, error)
	
	// Statistics
	CountByCategory(ctx context.Context) (map[string]int, error)
	CountUserInterests(ctx context.Context, userID string) (int, error)
}

