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

// LikeRepository defines like operations
type LikeRepository interface {
	// Like adds a like to a postfeed
	Like(ctx context.Context, postfeedID, userID string) error
	// Unlike removes a like from a postfeed
	Unlike(ctx context.Context, postfeedID, userID string) error
	// IsLiked checks if user has liked the postfeed
	IsLiked(ctx context.Context, postfeedID, userID string) (bool, error)
	// GetLikeCount returns the like count
	GetLikeCount(ctx context.Context, postfeedID string) (int, error)
}

// CommentRepository defines comment operations
type CommentRepository interface {
	// Create adds a comment
	Create(ctx context.Context, comment *Comment) error
	// GetByID retrieves a comment
	GetByID(ctx context.Context, id string) (*Comment, error)
	// Update updates a comment
	Update(ctx context.Context, id string, content string) error
	// Delete soft deletes a comment
	Delete(ctx context.Context, id string) error
	// ListByPostfeed retrieves comments for a postfeed
	ListByPostfeed(ctx context.Context, postfeedID string, limit, offset int) ([]Comment, int, error)
	// GetReplies retrieves replies to a comment
	GetReplies(ctx context.Context, parentID string) ([]Comment, error)
}

// ShareRepository defines share operations
type ShareRepository interface {
	// Create adds a share
	Create(ctx context.Context, share *Share) error
	// GetByID retrieves a share
	GetByID(ctx context.Context, id string) (*Share, error)
	// ListByPostfeed retrieves shares for a postfeed
	ListByPostfeed(ctx context.Context, postfeedID string, limit, offset int) ([]Share, int, error)
	// GetShareCount returns the share count
	GetShareCount(ctx context.Context, postfeedID string) (int, error)
}
