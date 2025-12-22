package education

// GET /v1/users/:user_id/education - List user's education
type ListEducationResponse struct {
	Education []Education `json:"education"`
	Total     int         `json:"total"`
}

// GET /v1/education/:id - Get education by ID
type GetEducationResponse struct {
	Education Education `json:"education"`
}

// POST /v1/users/:user_id/education - Add education entry
type CreateEducationRequest struct {
	School       string `json:"school" binding:"required,min=2,max=255"`
	Degree       string `json:"degree,omitempty"`
	FieldOfStudy string `json:"field_of_study,omitempty"`
	Location     string `json:"location,omitempty"`
	StartYear    *int   `json:"start_year,omitempty"`
	EndYear      *int   `json:"end_year,omitempty"`
	Current      bool   `json:"current"`
	LogoURL      string `json:"logo_url,omitempty"`
	Grade        string `json:"grade,omitempty"`
	Activities   string `json:"activities,omitempty"`
	Description  string `json:"description,omitempty"`
}

type CreateEducationResponse struct {
	Message   string    `json:"message"`
	Education Education `json:"education"`
}

// PUT /v1/education/:id - Update education entry
type UpdateEducationRequest struct {
	School       string `json:"school,omitempty"`
	Degree       string `json:"degree,omitempty"`
	FieldOfStudy string `json:"field_of_study,omitempty"`
	Location     string `json:"location,omitempty"`
	StartYear    *int   `json:"start_year,omitempty"`
	EndYear      *int   `json:"end_year,omitempty"`
	Current      *bool  `json:"current,omitempty"`
	LogoURL      string `json:"logo_url,omitempty"`
	Grade        string `json:"grade,omitempty"`
	Activities   string `json:"activities,omitempty"`
	Description  string `json:"description,omitempty"`
}

type UpdateEducationResponse struct {
	Message   string    `json:"message"`
	Education Education `json:"education"`
}

// DELETE /v1/education/:id - Delete education entry
type DeleteEducationResponse struct {
	Message string `json:"message"`
}
