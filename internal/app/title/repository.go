package title

import "context"

// TitleRepository defines methods for title data access
type TitleRepository interface {
	// Title CRUD
	Create(ctx context.Context, title *Title) error
	FindAll(ctx context.Context) ([]Title, error)
	FindByID(ctx context.Context, id string) (*Title, error)
	FindByName(ctx context.Context, name string) (*Title, error)
	FindByDefinedBy(ctx context.Context, definedBy string) ([]Title, error)
	Update(ctx context.Context, title *Title) error
	Delete(ctx context.Context, id string) error

	// Usage tracking
	IncrementUsedByCount(ctx context.Context, id string) (int, error)
	DecrementUsedByCount(ctx context.Context, id string) (int, error)

	// Statistics
	CountByDefinedBy(ctx context.Context) (map[string]int, error)
	GetMostUsedTitles(ctx context.Context, limit int) ([]Title, error)
}
