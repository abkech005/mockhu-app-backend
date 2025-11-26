package follow

import (
	"context"
	"errors"
	"strings"

	"mockhu-app-backend/internal/app/auth"
)

// Errors
var (
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
	ErrUserNotFound     = errors.New("user not found")
)

// FollowService defines the business logic for follow operations
type FollowService interface {
	Follow(ctx context.Context, followerID, followingID string) (*FollowResponse, error)
	Unfollow(ctx context.Context, followerID, followingID string) (*FollowResponse, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (*IsFollowingResponse, error)
	GetFollowers(ctx context.Context, userID, currentUserID string, page, limit int) (*UserListResponse, error)
	GetFollowing(ctx context.Context, userID, currentUserID string, page, limit int) (*UserListResponse, error)
	GetFollowStats(ctx context.Context, userID string) (*FollowStatsResponse, error)
}

// followService implements FollowService
type followService struct {
	followRepo FollowRepository
	userRepo   auth.UserRepository
}

// NewService creates a new follow service
func NewService(followRepo FollowRepository, userRepo auth.UserRepository) FollowService {
	return &followService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

// Follow creates a follow relationship
func (s *followService) Follow(ctx context.Context, followerID, followingID string) (*FollowResponse, error) {
	// Check if trying to follow self
	if followerID == followingID {
		return nil, ErrCannotFollowSelf
	}

	// Check if target user exists
	user, err := s.userRepo.FindByID(ctx, followingID)
	if err != nil {
		// Check if error indicates user not found
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Create follow relationship
	err = s.followRepo.Follow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		Message:     "followed successfully",
		IsFollowing: true,
	}, nil
}

// Unfollow removes a follow relationship
func (s *followService) Unfollow(ctx context.Context, followerID, followingID string) (*FollowResponse, error) {
	// Check if trying to unfollow self
	if followerID == followingID {
		return nil, ErrCannotFollowSelf
	}

	// Remove follow relationship
	err := s.followRepo.Unfollow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		Message:     "unfollowed successfully",
		IsFollowing: false,
	}, nil
}

// IsFollowing checks if a user follows another user
func (s *followService) IsFollowing(ctx context.Context, followerID, followingID string) (*IsFollowingResponse, error) {
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &IsFollowingResponse{
		IsFollowing: isFollowing,
	}, nil
}

// GetFollowers returns list of users who follow the given user
func (s *followService) GetFollowers(ctx context.Context, userID, currentUserID string, page, limit int) (*UserListResponse, error) {
	offset := (page - 1) * limit

	// Get followers
	follows, err := s.followRepo.GetFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := s.followRepo.GetFollowerCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build user list with details
	users := make([]UserListItem, 0, len(follows))
	for _, f := range follows {
		// Get user details
		user, err := s.userRepo.FindByID(ctx, f.FollowerID)
		if err != nil || user == nil {
			continue
		}

		// Check if current user follows this follower
		isFollowedByMe := false
		if currentUserID != "" && currentUserID != f.FollowerID {
			isFollowedByMe, _ = s.followRepo.IsFollowing(ctx, currentUserID, f.FollowerID)
		}

		users = append(users, UserListItem{
			ID:             user.ID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			AvatarURL:      user.AvatarURL,
			IsFollowedByMe: isFollowedByMe,
		})
	}

	return &UserListResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetFollowing returns list of users that the given user follows
func (s *followService) GetFollowing(ctx context.Context, userID, currentUserID string, page, limit int) (*UserListResponse, error) {
	offset := (page - 1) * limit

	// Get following
	follows, err := s.followRepo.GetFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := s.followRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build user list with details
	users := make([]UserListItem, 0, len(follows))
	for _, f := range follows {
		// Get user details
		user, err := s.userRepo.FindByID(ctx, f.FollowingID)
		if err != nil || user == nil {
			continue
		}

		// Check if current user follows this user
		isFollowedByMe := false
		if currentUserID != "" && currentUserID != f.FollowingID {
			isFollowedByMe, _ = s.followRepo.IsFollowing(ctx, currentUserID, f.FollowingID)
		}

		users = append(users, UserListItem{
			ID:             user.ID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			AvatarURL:      user.AvatarURL,
			IsFollowedByMe: isFollowedByMe,
		})
	}

	return &UserListResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetFollowStats returns follower/following counts for a user
func (s *followService) GetFollowStats(ctx context.Context, userID string) (*FollowStatsResponse, error) {
	stats, err := s.followRepo.GetFollowStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &FollowStatsResponse{
		UserID:         userID,
		FollowerCount:  stats.FollowerCount,
		FollowingCount: stats.FollowingCount,
	}, nil
}
