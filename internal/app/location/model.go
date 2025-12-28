package location

import "time"

// Location represents a city/country location
type Location struct {
	ID          string    `json:"id"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	UsedByCount int       `json:"used_by_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
