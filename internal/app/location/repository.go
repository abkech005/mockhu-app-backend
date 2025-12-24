package location

import "context"

// LocationRepository defines methods for location data access
type LocationRepository interface {
	// CRUD
	Create(ctx context.Context, location *Location) error
	FindByID(ctx context.Context, id string) (*Location, error)
	FindAll(ctx context.Context) ([]Location, error)
	Update(ctx context.Context, location *Location) error
	Delete(ctx context.Context, id string) error

	// Search
	Search(ctx context.Context, query string) ([]Location, error)
	FindByCity(ctx context.Context, city string) ([]Location, error)
	FindByCountry(ctx context.Context, country string) ([]Location, error)
	FindByCityAndCountry(ctx context.Context, city, country string) (*Location, error)

	// Usage tracking
	IncrementUsedByCount(ctx context.Context, id string) (int, error)
	DecrementUsedByCount(ctx context.Context, id string) (int, error)

	// Statistics
	GetMostUsedLocations(ctx context.Context, limit int) ([]Location, error)
	CountByCountry(ctx context.Context) (map[string]int, error)
}
