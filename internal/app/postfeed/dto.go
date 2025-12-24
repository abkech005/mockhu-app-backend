package postfeed

import "encoding/json"

// --- Request DTOs ---

// CreatePostfeedRequest for creating any type of postfeed
type CreatePostfeedRequest struct {
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Content     string          `json:"content,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	IsAnonymous bool            `json:"is_anonymous,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

// UpdatePostfeedRequest for updating a postfeed
type UpdatePostfeedRequest struct {
	Title    string          `json:"title,omitempty"`
	Content  string          `json:"content,omitempty"`
	Tags     []string        `json:"tags,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// ListPostfeedsRequest for filtering postfeeds
type ListPostfeedsRequest struct {
	Type   string `query:"type"`
	UserID string `query:"user_id"`
	Tag    string `query:"tag"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}

// --- Response DTOs ---

// PostfeedResponse for a single postfeed with author info
type PostfeedResponse struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	Author       *AuthorInfo     `json:"author,omitempty"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Content      string          `json:"content,omitempty"`
	Tags         []string        `json:"tags,omitempty"`
	IsAnonymous  bool            `json:"is_anonymous"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	ViewCount    int             `json:"view_count"`
	LikeCount    int             `json:"like_count"`
	CommentCount int             `json:"comment_count"`
	ShareCount   int             `json:"share_count"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// AuthorInfo for non-anonymous posts
type AuthorInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ListPostfeedsResponse for paginated list
type ListPostfeedsResponse struct {
	Postfeeds  []PostfeedResponse `json:"postfeeds"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// CreatePostfeedResponse after creation
type CreatePostfeedResponse struct {
	Message  string           `json:"message"`
	Postfeed PostfeedResponse `json:"postfeed"`
}

// UpdatePostfeedResponse after update
type UpdatePostfeedResponse struct {
	Message  string           `json:"message"`
	Postfeed PostfeedResponse `json:"postfeed"`
}

// DeletePostfeedResponse after deletion
type DeletePostfeedResponse struct {
	Message string `json:"message"`
}
