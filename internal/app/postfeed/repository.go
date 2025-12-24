package postfeed

import "context"

// Repository defines the interface for postfeed data operations
type Repository interface {
	// Create a new postfeed
	Create(ctx context.Context, postfeed *Postfeed) error

	// GetByID retrieves a postfeed by ID
	GetByID(ctx context.Context, id string) (*Postfeed, error)

	// Update a postfeed
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete (soft delete) a postfeed
	Delete(ctx context.Context, id string) error

	// List postfeeds with filters and pagination
	List(ctx context.Context, filter ListFilter) ([]Postfeed, int, error)

	// GetByUserID retrieves all postfeeds by a user
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]Postfeed, int, error)

	// GetByType retrieves postfeeds by type
	GetByType(ctx context.Context, postType string, limit, offset int) ([]Postfeed, int, error)

	// IncrementViewCount increments view count
	IncrementViewCount(ctx context.Context, id string) error

	// UpdateEngagementCounts updates like/comment/share counts
	UpdateEngagementCounts(ctx context.Context, id string, likes, comments, shares int) error
}

// ListFilter for filtering postfeeds
type ListFilter struct {
	Type   string
	UserID string
	Tag    string
	Limit  int
	Offset int
}
