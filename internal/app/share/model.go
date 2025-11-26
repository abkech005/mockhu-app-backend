package share

import "time"

// Share represents a post share
type Share struct {
	ID           string    `json:"id"`
	PostID       string    `json:"post_id"`
	UserID       string    `json:"user_id"`
	SharedToType string    `json:"shared_to_type"`
	CreatedAt    time.Time `json:"created_at"`
}


