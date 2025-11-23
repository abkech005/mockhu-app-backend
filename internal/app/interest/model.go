package interest

import "time"

// Interest represents a predefined interest that users can select
type Interest struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Category  string    `json:"category"`
	Icon      string    `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInterest represents a user's selected interest
type UserInterest struct {
	UserID     string    `json:"user_id"`
	InterestID string    `json:"interest_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Interest categories
const (
	CategoryTechnology    = "technology"
	CategoryArts          = "arts"
	CategorySports        = "sports"
	CategoryEntertainment = "entertainment"
	CategoryLifestyle     = "lifestyle"
	CategoryBusiness      = "business"
	CategoryEducation     = "education"
	CategorySocial        = "social"
)

// UserInterestDetail includes interest details for a user
type UserInterestDetail struct {
	UserID    string    `json:"user_id"`
	Interest  Interest  `json:"interest"`
	CreatedAt time.Time `json:"created_at"`
}
