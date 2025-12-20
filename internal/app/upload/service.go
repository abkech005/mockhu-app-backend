package upload

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/infra/r2"
)

const (
	maxFileSize   = 5 * 1024 * 1024 // 5MB
	presignExpiry = 5 * time.Minute
)

var allowedContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type Service struct {
	r2Client *r2.Client
	userRepo auth.UserRepository
}

func NewService(r2Client *r2.Client, userRepo auth.UserRepository) *Service {
	return &Service{
		r2Client: r2Client,
		userRepo: userRepo,
	}
}

// RequestAvatarUpload generates a presigned PUT URL for direct upload.
func (s *Service) RequestAvatarUpload(ctx context.Context, req *AvatarUploadRequest) (*AvatarUploadResponse, error) {
	// Validate content type
	ext, ok := allowedContentTypes[req.ContentType]
	if !ok {
		return nil, fmt.Errorf("invalid content type: must be image/jpeg, image/png, or image/webp")
	}

	// Validate file size
	if req.FileSize > maxFileSize {
		return nil, fmt.Errorf("file too large: max 5MB allowed")
	}

	// Generate unique key: avatars/{user_id}/{timestamp}{ext}
	fileKey := fmt.Sprintf("avatars/%s/%d%s", req.UserID, time.Now().UnixNano(), ext)

	// Generate presigned URL
	uploadURL, err := s.r2Client.GeneratePresignedPutURL(ctx, fileKey, req.ContentType, presignExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &AvatarUploadResponse{
		UploadURL: uploadURL,
		FileKey:   fileKey,
		ExpiresIn: int(presignExpiry.Seconds()),
	}, nil
}

// ConfirmAvatarUpload saves the avatar URL to the user profile.
func (s *Service) ConfirmAvatarUpload(ctx context.Context, req *AvatarConfirmRequest) (*AvatarConfirmResponse, error) {
	// Validate file key belongs to user
	if !strings.HasPrefix(req.FileKey, fmt.Sprintf("avatars/%s/", req.UserID)) {
		return nil, fmt.Errorf("invalid file key")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(req.FileKey))
	if ext != ".jpg" && ext != ".png" && ext != ".webp" {
		return nil, fmt.Errorf("invalid file extension")
	}

	// Get public URL
	avatarURL := s.r2Client.GetPublicURL(req.FileKey)

	// Update user profile
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Delete old avatar if exists
	if user.AvatarURL != "" && strings.Contains(user.AvatarURL, "avatar.mockhu.com") {
		oldKey := strings.TrimPrefix(user.AvatarURL, s.r2Client.GetPublicURL("")+"/")
		if oldKey != "" {
			_ = s.r2Client.DeleteObject(ctx, oldKey) // Best effort cleanup
		}
	}

	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &AvatarConfirmResponse{
		AvatarURL: avatarURL,
	}, nil
}
