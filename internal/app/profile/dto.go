package profile

// ProfileStats represents user statistics
type ProfileStats struct {
	PostsCount     int `json:"posts_count"`
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
}

// InstitutionInfo represents basic institution information
type InstitutionInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProfileResponse for public profile view (GET /v1/users/:userId/profile)
type ProfileResponse struct {
	ID                     string           `json:"id"`
	Username               string           `json:"username"`
	FirstName              string           `json:"first_name"`
	LastName               string           `json:"last_name"`
	AvatarURL              string           `json:"avatar_url,omitempty"`
	Bio                    string           `json:"bio,omitempty"`
	Institution            *InstitutionInfo `json:"institution,omitempty"`
	Stats                  ProfileStats     `json:"stats"`
	IsFollowing            bool             `json:"is_following"`
	IsFollowedBy           bool             `json:"is_followed_by"`
	MutualConnectionsCount int              `json:"mutual_connections_count"`
	CreatedAt              string           `json:"created_at"`
}

// PrivacySettings represents user privacy settings
type PrivacySettings struct {
	WhoCanMessage     string `json:"who_can_message"`
	WhoCanSeePosts    string `json:"who_can_see_posts"`
	ShowFollowersList bool   `json:"show_followers_list"`
	ShowFollowingList bool   `json:"show_following_list"`
}

// OwnProfileResponse for authenticated user's own profile (GET /v1/users/me/profile)
type OwnProfileResponse struct {
	ID                  string           `json:"id"`
	Username            string           `json:"username"`
	FirstName           string           `json:"first_name"`
	LastName            string           `json:"last_name"`
	Email               string           `json:"email"`
	Phone               string           `json:"phone"`
	DateOfBirth         string           `json:"date_of_birth"`
	AvatarURL           string           `json:"avatar_url,omitempty"`
	Bio                 string           `json:"bio,omitempty"`
	Institution         *InstitutionInfo `json:"institution,omitempty"`
	Stats               ProfileStats     `json:"stats"`
	PrivacySettings     PrivacySettings  `json:"privacy_settings"`
	EmailVerified       bool             `json:"email_verified"`
	PhoneVerified       bool             `json:"phone_verified"`
	OnboardingCompleted bool             `json:"onboarding_completed"`
	CreatedAt           string           `json:"created_at"`
}

// UpdateProfileRequest for updating user profile (PUT /v1/users/me/profile)
type UpdateProfileRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

// UpdatePrivacyRequest for updating privacy settings (PUT /v1/users/me/privacy)
type UpdatePrivacyRequest struct {
	WhoCanMessage     string `json:"who_can_message,omitempty"`
	WhoCanSeePosts    string `json:"who_can_see_posts,omitempty"`
	ShowFollowersList *bool  `json:"show_followers_list,omitempty"`
	ShowFollowingList *bool  `json:"show_following_list,omitempty"`
}

// MutualConnectionUser represents a user in mutual connections list
type MutualConnectionUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// MutualConnectionsResponse for mutual connections list (GET /v1/users/:userId/mutual-connections)
type MutualConnectionsResponse struct {
	MutualConnections []*MutualConnectionUser `json:"mutual_connections"`
	Pagination        PaginationInfo          `json:"pagination"`
}

// MutualConnectionsCountResponse for mutual connections count
type MutualConnectionsCountResponse struct {
	UserID                 string `json:"user_id"`
	MutualConnectionsCount int    `json:"mutual_connections_count"`
}

// PaginationInfo for paginated responses
type PaginationInfo struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Limit      int `json:"limit"`
}

// AvatarUploadResponse for avatar upload (POST /v1/users/me/avatar)
type AvatarUploadResponse struct {
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}
