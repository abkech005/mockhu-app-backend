package post

import "context"

// PostRepository defines the interface for post data operations
type PostRepository interface {
	// Post CRUD operations
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id string) (*Post, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error

	// Feed operations
	GetFeed(ctx context.Context, userID string, limit, offset int) ([]*Post, error)

	// Reaction operations
	AddReaction(ctx context.Context, reaction *Reaction) error
	RemoveReaction(ctx context.Context, postID, userID string) error
	GetReactions(ctx context.Context, postID string, limit, offset int) ([]*Reaction, error)
	GetReactionCount(ctx context.Context, postID string) (int, error)
	HasUserReacted(ctx context.Context, postID, userID string) (bool, error)
}
