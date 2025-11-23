package interest

// GET /v1/interests - List all interests
type ListInterestsResponse struct {
	Interests []Interest `json:"interests"`
	Total     int        `json:"total"`
}

// GET /v1/interests/categories - List categories
type ListCategoriesResponse struct {
	Categories []CategoryInfo `json:"categories"`
}

type CategoryInfo struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GET /v1/users/:id/interests - Get user's interests
type GetUserInterestsResponse struct {
	UserID    string     `json:"user_id"`
	Interests []Interest `json:"interests"`
	Count     int        `json:"count"`
}

// POST /v1/users/:id/interests - Add interests to user
type AddUserInterestsRequest struct {
	InterestSlugs []string `json:"interest_slugs" binding:"required,min=1"`
}

type AddUserInterestsResponse struct {
	Message        string     `json:"message"`
	AddedInterests []Interest `json:"added_interests"`
	TotalCount     int        `json:"total_count"`
}

// DELETE /v1/users/:id/interests/:interest_slug - Remove interest
type RemoveUserInterestResponse struct {
	Message string `json:"message"`
}

// PUT /v1/users/:id/interests - Replace all user interests
type ReplaceUserInterestsRequest struct {
	InterestSlugs []string `json:"interest_slugs" binding:"required"`
}

type ReplaceUserInterestsResponse struct {
	Message   string     `json:"message"`
	Interests []Interest `json:"interests"`
	Count     int        `json:"count"`
}

// POST /v1/interests - Create new interest (admin only)
type CreateInterestRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Slug     string `json:"slug" binding:"required,min=2,max=100"`
	Category string `json:"category" binding:"required"`
	Icon     string `json:"icon,omitempty"`
}

type CreateInterestResponse struct {
	Message  string   `json:"message"`
	Interest Interest `json:"interest"`
}
