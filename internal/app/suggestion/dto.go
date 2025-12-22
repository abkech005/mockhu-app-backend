package suggestion

// SuggestedUser represents a user suggestion with metadata
type SuggestedUser struct {
	UserID               string   `json:"user_id"`
	Username             string   `json:"username"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	AvatarURL            string   `json:"avatar_url"`
	Bio                  string   `json:"bio"`
	Place                string   `json:"place"`
	SharedInterests      []string `json:"shared_interests"`
	SharedInterestsCount int      `json:"shared_interests_count"`
	FollowerCount        int      `json:"follower_count"`
	Reason               string   `json:"reason"`
}

// GET /v1/suggestions/users
type GetUserSuggestionsResponse struct {
	Suggestions []SuggestedUser `json:"suggestions"`
	Total       int             `json:"total"`
}
