package onboarding

// POST /v1/onboard/basic
type BasicRequest struct {
	TempID    string `json:"temp_id,omitempty"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	DOB       string `json:"dob" binding:"required"`
}

type BasicResponse struct {
	TempID  string `json:"temp_id"`
	Message string `json:"message"`
}

// POST /v1/onboard/profile
type ProfileRequest struct {
	TempID       string        `json:"temp_id" binding:"required"`
	Username     string        `json:"username" binding:"required"`
	AvatarUpload *AvatarUpload `json:"avatar_upload,omitempty"`
}

type AvatarUpload struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type ProfileResponse struct {
	TempID            string `json:"temp_id"`
	UsernameAvailable bool   `json:"username_available"`
	AvatarURL         string `json:"avatar_url,omitempty"`
}

// POST /v1/onboard/interests
type InterestsRequest struct {
	Interests []string `json:"interests" binding:"required"`
}

type InterestsResponse struct {
	Message string `json:"message"`
}
