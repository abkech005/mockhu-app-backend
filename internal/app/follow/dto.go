package follow

// FollowResponse is returned after follow/unfollow action
type FollowResponse struct {
	Message     string `json:"message"`
	IsFollowing bool   `json:"is_following"`
}

// FollowStatsResponse contains user follow statistics
type FollowStatsResponse struct {
	UserID         string `json:"user_id"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
}

// UserListItem represents a user in follower/following lists
type UserListItem struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	AvatarURL      string `json:"avatar_url"`
	IsFollowedByMe bool   `json:"is_followed_by_me"` // For "Follow back" button
}

// UserListResponse is the response for follower/following lists
type UserListResponse struct {
	Users      []UserListItem `json:"users"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}

// IsFollowingResponse for checking follow status
type IsFollowingResponse struct {
	IsFollowing bool `json:"is_following"`
}
