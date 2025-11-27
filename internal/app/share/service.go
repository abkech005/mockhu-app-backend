package share

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/post"
)

// Errors
var (
	ErrShareNotFound    = errors.New("share not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrPostNotFound     = errors.New("post not found")
	ErrInvalidShareType = errors.New("invalid share type")
	ErrAlreadyShared    = errors.New("post already shared by user")
)

// ShareService defines the business logic for share operations
type ShareService interface {
	CreateShare(ctx context.Context, postID, userID string, req *CreateShareRequest) (*ShareResponse, error)
	GetShare(ctx context.Context, shareID, currentUserID string) (*ShareResponse, error)
	GetPostShares(ctx context.Context, postID, currentUserID string, page, limit int) (*ShareListResponse, error)
	GetUserShares(ctx context.Context, userID, currentUserID string, page, limit int) (*ShareListResponse, error)
	DeleteShare(ctx context.Context, shareID, userID string) error
	GetShareCount(ctx context.Context, postID string) (int, error)
	HasUserShared(ctx context.Context, postID, userID string) (bool, error)
}

// shareService implements ShareService
type shareService struct {
	shareRepo ShareRepository
	userRepo  auth.UserRepository
	postRepo  post.PostRepository
}

// NewService creates a new share service
func NewService(shareRepo ShareRepository, userRepo auth.UserRepository, postRepo post.PostRepository) ShareService {
	return &shareService{
		shareRepo: shareRepo,
		userRepo:  userRepo,
		postRepo:  postRepo,
	}
}

// CreateShare creates a new share for a post
func (s *shareService) CreateShare(ctx context.Context, postID, userID string, req *CreateShareRequest) (*ShareResponse, error) {
	// Validate share type
	if req.SharedToType != "timeline" && req.SharedToType != "dm" && req.SharedToType != "external" {
		return nil, ErrInvalidShareType
	}

	// Validate post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Check if user has already shared this post (prevent duplicate shares)
	hasShared, err := s.shareRepo.HasUserShared(ctx, postID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user shared: %w", err)
	}
	if hasShared {
		return nil, ErrAlreadyShared
	}

	// Create share
	share := &Share{
		PostID:       postID,
		UserID:       userID,
		SharedToType: req.SharedToType,
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, fmt.Errorf("failed to create share: %w", err)
	}

	// Get user info
	user, err := s.getUserInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Build response
	response := &ShareResponse{
		ID:           share.ID,
		PostID:       share.PostID,
		User:         *user,
		SharedToType: share.SharedToType,
		CreatedAt:    share.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// GetShare retrieves a single share by ID
func (s *shareService) GetShare(ctx context.Context, shareID, currentUserID string) (*ShareResponse, error) {
	share, err := s.shareRepo.GetByID(ctx, shareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get share: %w", err)
	}
	if share == nil {
		return nil, ErrShareNotFound
	}

	return s.convertToResponse(ctx, share)
}

// GetPostShares retrieves all shares for a post
func (s *shareService) GetPostShares(ctx context.Context, postID, currentUserID string, page, limit int) (*ShareListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get shares
	shares, err := s.shareRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get shares: %w", err)
	}

	// Convert to response
	shareResponses := make([]*ShareResponse, 0, len(shares))
	for _, share := range shares {
		response, err := s.convertToResponse(ctx, share)
		if err != nil {
			continue // Skip shares with errors
		}
		shareResponses = append(shareResponses, response)
	}

	// Get total count for pagination
	totalCount, _ := s.shareRepo.GetShareCount(ctx, postID)
	totalPages := (totalCount + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &ShareListResponse{
		Shares: shareResponses,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: totalCount,
			Limit:      limit,
		},
	}, nil
}

// GetUserShares retrieves all shares by a user
func (s *shareService) GetUserShares(ctx context.Context, userID, currentUserID string, page, limit int) (*ShareListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get shares
	shares, err := s.shareRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get shares: %w", err)
	}

	// Convert to response
	shareResponses := make([]*ShareResponse, 0, len(shares))
	for _, share := range shares {
		response, err := s.convertToResponse(ctx, share)
		if err != nil {
			continue // Skip shares with errors
		}
		shareResponses = append(shareResponses, response)
	}

	// Calculate total pages (simplified)
	totalPages := 1
	if len(shareResponses) == limit {
		totalPages = page + 1
	}

	return &ShareListResponse{
		Shares: shareResponses,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: len(shareResponses),
			Limit:      limit,
		},
	}, nil
}

// DeleteShare deletes a share
func (s *shareService) DeleteShare(ctx context.Context, shareID, userID string) error {
	// Get share to verify ownership
	share, err := s.shareRepo.GetByID(ctx, shareID)
	if err != nil {
		return fmt.Errorf("failed to get share: %w", err)
	}
	if share == nil {
		return ErrShareNotFound
	}

	// Check ownership
	if share.UserID != userID {
		return ErrUnauthorized
	}

	// Delete share
	err = s.shareRepo.Delete(ctx, shareID)
	if err != nil {
		return fmt.Errorf("failed to delete share: %w", err)
	}

	return nil
}

// GetShareCount returns the total number of shares for a post
func (s *shareService) GetShareCount(ctx context.Context, postID string) (int, error) {
	return s.shareRepo.GetShareCount(ctx, postID)
}

// HasUserShared checks if a user has shared a post
func (s *shareService) HasUserShared(ctx context.Context, postID, userID string) (bool, error) {
	return s.shareRepo.HasUserShared(ctx, postID, userID)
}

// Helper methods

// getUserInfo retrieves user information for a share
func (s *shareService) getUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		AvatarURL: user.AvatarURL,
	}, nil
}

// convertToResponse converts a share to a response
func (s *shareService) convertToResponse(ctx context.Context, share *Share) (*ShareResponse, error) {
	// Get user info
	user, err := s.getUserInfo(ctx, share.UserID)
	if err != nil {
		return nil, err
	}

	return &ShareResponse{
		ID:           share.ID,
		PostID:       share.PostID,
		User:         *user,
		SharedToType: share.SharedToType,
		CreatedAt:    share.CreatedAt.Format(time.RFC3339),
	}, nil
}

