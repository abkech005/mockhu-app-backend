package post

// CreatePostRequest is the request DTO for creating a new post
type CreatePostRequest struct {
	Content     string   `json:"content" validate:"required,min=1,max=5000"`
	Images      []string `json:"images" validate:"max=10"`
	IsAnonymous bool     `json:"is_anonymous"`
}

// PostResponse is the response DTO for a post with enriched data
type PostResponse struct {
	ID        string       `json:"id"`
	Author    AuthorInfo   `json:"author"`
	Content   string       `json:"content"`
	Images    []string     `json:"images"`
	Reactions ReactionInfo `json:"reactions"`
	CreatedAt string       `json:"created_at"`
}

// AuthorInfo contains author information for a post
type AuthorInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	AvatarURL string `json:"avatar_url"`
}

// ReactionInfo contains reaction information for a post
type ReactionInfo struct {
	FireCount   int          `json:"fire_count"`
	IsFiredByMe bool         `json:"is_fired_by_me"`
	RecentUsers []AuthorInfo `json:"recent_users"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Limit      int `json:"limit"`
}

// FeedResponse is the response for feed endpoint
type FeedResponse struct {
	Posts      []*PostResponse `json:"posts"`
	Pagination PaginationInfo  `json:"pagination"`
}

// ReactionResponse is the response for toggling a reaction
type ReactionResponse struct {
	PostID      string `json:"post_id"`
	FireCount   int    `json:"fire_count"`
	IsFiredByMe bool   `json:"is_fired_by_me"`
}

