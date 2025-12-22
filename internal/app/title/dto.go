package title

// GET /v1/titles - List all titles
type ListTitlesResponse struct {
	Titles []Title `json:"titles"`
	Total  int     `json:"total"`
}

// GET /v1/titles/:id - Get title by ID
type GetTitleResponse struct {
	Title Title `json:"title"`
}

// POST /v1/titles - Create new title (admin only)
type CreateTitleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type CreateTitleResponse struct {
	Message string `json:"message"`
	Title   Title  `json:"title"`
}

// PUT /v1/titles/:id - Update title (admin only)
type UpdateTitleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateTitleResponse struct {
	Message string `json:"message"`
	Title   Title  `json:"title"`
}

// DELETE /v1/titles/:id - Delete title (admin only)
type DeleteTitleResponse struct {
	Message string `json:"message"`
}

// POST /v1/titles/:id/increment - Increment used_by_count
type IncrementUsageResponse struct {
	Message     string `json:"message"`
	UsedByCount int    `json:"used_by_count"`
}

// POST /v1/titles/:id/decrement - Decrement used_by_count
type DecrementUsageResponse struct {
	Message     string `json:"message"`
	UsedByCount int    `json:"used_by_count"`
}
