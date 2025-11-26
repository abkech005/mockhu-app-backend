package post

import "time"

// Post represents a user's post
type Post struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	Images      []string  `json:"images"`
	IsAnonymous bool      `json:"is_anonymous"`
	IsActive    bool      `json:"is_active"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Reaction represents a user's reaction to a post
type Reaction struct {
	ID           string    `json:"id"`
	PostID       string    `json:"post_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}

