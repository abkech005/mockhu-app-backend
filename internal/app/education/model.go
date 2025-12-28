package education

import "time"

// Education represents a user's education entry
type Education struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	School       string    `json:"school"`
	Degree       string    `json:"degree,omitempty"`
	FieldOfStudy string    `json:"field_of_study,omitempty"`
	Location     string    `json:"location,omitempty"`
	StartYear    *int      `json:"start_year,omitempty"`
	EndYear      *int      `json:"end_year,omitempty"`
	Current      bool      `json:"current"`
	LogoURL      string    `json:"logo_url,omitempty"`
	Grade        string    `json:"grade,omitempty"`
	Activities   string    `json:"activities,omitempty"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
