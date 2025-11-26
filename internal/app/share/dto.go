package share

// CreateShareRequest is the request DTO for creating a share
type CreateShareRequest struct {
	SharedToType string `json:"shared_to_type" validate:"required,oneof=timeline dm external"`
}

// ShareResponse is the response DTO for a share with enriched data
type ShareResponse struct {
	ID           string   `json:"id"`
	PostID       string   `json:"post_id"`
	User         UserInfo `json:"user"`
	SharedToType string   `json:"shared_to_type"`
	CreatedAt    string   `json:"created_at"`
}

// UserInfo contains user information for a share
type UserInfo struct {
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

// ShareListResponse is the response for listing shares
type ShareListResponse struct {
	Shares     []*ShareResponse `json:"shares"`
	Pagination PaginationInfo    `json:"pagination"`
}


