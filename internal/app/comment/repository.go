package comment

import "context"

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	// Comment CRUD operations
	Create(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, id string) (*Comment, error)
	GetByPostID(ctx context.Context, postID string, limit, offset int) ([]*Comment, error)
	GetReplies(ctx context.Context, parentCommentID string, limit, offset int) ([]*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id string) error

	// Count operations
	GetCommentCount(ctx context.Context, postID string) (int, error)
	GetReplyCount(ctx context.Context, parentCommentID string) (int, error)
}

