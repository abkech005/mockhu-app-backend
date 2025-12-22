package title

import "time"

// Title represents a user title (e.g., Student, Teacher, CAT Aspirant)
type Title struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DefinedBy   string    `json:"defined_by"` // "admin" or "user"
	UsedByCount int       `json:"used_by_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DefinedBy constants
const (
	DefinedByAdmin = "admin"
	DefinedByUser  = "user"
)
