package comment

import "time"

// Comment represents a comment on a post
type Comment struct {
	ID              string    `json:"id"`
	PostID          string    `json:"post_id"`
	UserID          string    `json:"user_id"`
	ParentCommentID *string   `json:"parent_comment_id,omitempty"`
	Content         string    `json:"content"`
	IsAnonymous     bool      `json:"is_anonymous"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

