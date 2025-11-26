package comment

// CreateCommentRequest is the request DTO for creating a comment
type CreateCommentRequest struct {
	Content         string  `json:"content" validate:"required,min=1,max=2000"`
	ParentCommentID *string `json:"parent_comment_id,omitempty" validate:"omitempty,uuid"`
	IsAnonymous     bool    `json:"is_anonymous"`
}

// CommentResponse is the response DTO for a comment with enriched data
type CommentResponse struct {
	ID              string       `json:"id"`
	PostID          string       `json:"post_id"`
	Author          AuthorInfo   `json:"author"`
	Content         string       `json:"content"`
	ParentCommentID *string      `json:"parent_comment_id,omitempty"`
	Replies         []*CommentResponse `json:"replies,omitempty"`
	ReplyCount      int          `json:"reply_count"`
	CreatedAt       string       `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
}

// AuthorInfo contains author information for a comment
type AuthorInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	AvatarURL string `json:"avatar_url"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Limit      int `json:"limit"`
}

// CommentListResponse is the response for listing comments
type CommentListResponse struct {
	Comments   []*CommentResponse `json:"comments"`
	Pagination PaginationInfo     `json:"pagination"`
}

