package upload

// POST /v1/upload/avatar/request
type AvatarUploadRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	ContentType string `json:"content_type"` // image/jpeg, image/png, image/webp
	FileSize    int64  `json:"file_size"`    // For validation (max 5MB)
}

type AvatarUploadResponse struct {
	UploadURL string `json:"upload_url"` // Presigned PUT URL
	FileKey   string `json:"file_key"`   // Key for confirm step
	ExpiresIn int    `json:"expires_in"` // Seconds
}

// POST /v1/upload/avatar/confirm
type AvatarConfirmRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	FileKey string `json:"file_key" binding:"required"`
}

type AvatarConfirmResponse struct {
	AvatarURL string `json:"avatar_url"` // Public CDN URL
}

// POST /v1/upload/media/request - For postfeed media attachments
type MediaUploadRequest struct {
	UserID      string `json:"user_id"`
	ContentType string `json:"content_type"` // image/*, video/*, audio/*, application/pdf
	FileSize    int64  `json:"file_size"`    // For validation
	FileName    string `json:"file_name"`    // Original filename
	MediaType   string `json:"media_type"`   // image, video, audio, document
}

type MediaUploadResponse struct {
	UploadURL    string `json:"upload_url"`    // Presigned PUT URL
	FileKey      string `json:"file_key"`      // R2 object key
	PublicURL    string `json:"public_url"`    // Public CDN URL (for use in postfeed)
	ThumbnailURL string `json:"thumbnail_url"` // Thumbnail URL (for images/videos)
	ExpiresIn    int    `json:"expires_in"`    // Seconds
}
