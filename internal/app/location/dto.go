package location

// GET /v1/locations - List locations
type ListLocationsResponse struct {
	Locations []Location `json:"locations"`
	Total     int        `json:"total"`
}

// GET /v1/locations/search - Search locations
type SearchLocationsResponse struct {
	Locations []Location `json:"locations"`
	Total     int        `json:"total"`
}

// GET /v1/locations/:id - Get location by ID
type GetLocationResponse struct {
	Location Location `json:"location"`
}

// POST /v1/locations - Create location
type CreateLocationRequest struct {
	City    string `json:"city" binding:"required,min=2,max=100"`
	Country string `json:"country" binding:"required,min=2,max=100"`
}

type CreateLocationResponse struct {
	Message  string   `json:"message"`
	Location Location `json:"location"`
}

// PUT /v1/locations/:id - Update location
type UpdateLocationRequest struct {
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
}

type UpdateLocationResponse struct {
	Message  string   `json:"message"`
	Location Location `json:"location"`
}

// DELETE /v1/locations/:id - Delete location
type DeleteLocationResponse struct {
	Message string `json:"message"`
}

// POST /v1/locations/:id/increment - Increment count
type IncrementUsageResponse struct {
	Message     string `json:"message"`
	UsedByCount int    `json:"used_by_count"`
}

// POST /v1/locations/:id/decrement - Decrement count
type DecrementUsageResponse struct {
	Message     string `json:"message"`
	UsedByCount int    `json:"used_by_count"`
}
