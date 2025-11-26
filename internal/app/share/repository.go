package share

import "context"

// ShareRepository defines the interface for share data operations
type ShareRepository interface {
	// Share CRUD operations
	Create(ctx context.Context, share *Share) error
	GetByID(ctx context.Context, id string) (*Share, error)
	GetByPostID(ctx context.Context, postID string, limit, offset int) ([]*Share, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Share, error)
	Delete(ctx context.Context, id string) error

	// Count operations
	GetShareCount(ctx context.Context, postID string) (int, error)
	HasUserShared(ctx context.Context, postID, userID string) (bool, error)
}


