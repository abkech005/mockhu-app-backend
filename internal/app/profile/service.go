package profile

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ProfileService defines the business logic for profile operations
type ProfileService interface {
	// Profile viewing
	GetUserProfile(ctx context.Context, userID, currentUserID string) (*ProfileResponse, error)
	GetOwnProfile(ctx context.Context, userID string) (*OwnProfileResponse, error)

	// Profile management
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error)

	// Avatar management
	UploadAvatar(ctx context.Context, userID string, fileBytes []byte, filename string) (*AvatarUploadResponse, error)
	DeleteAvatar(ctx context.Context, userID string) error

	// Privacy settings
	GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, userID string, req *UpdatePrivacyRequest) (*PrivacySettings, error)

	// Mutual connections
	GetMutualConnections(ctx context.Context, currentUserID, targetUserID string, page, limit int) (*MutualConnectionsResponse, error)
	GetMutualConnectionsCount(ctx context.Context, currentUserID, targetUserID string) (int, error)
}

// profileService implements ProfileService
type profileService struct {
	profileRepo ProfileRepository
	db          *pgxpool.Pool
}

// NewService creates a new profile service
func NewService(profileRepo ProfileRepository, db *pgxpool.Pool) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		db:          db,
	}
}

// GetUserProfile retrieves a public profile view
func (s *profileService) GetUserProfile(ctx context.Context, userID, currentUserID string) (*ProfileResponse, error) {
	// Get user from repository
	user, err := s.profileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get profile stats
	stats, err := s.getProfileStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile stats: %w", err)
	}

	// Build response
	response := &ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Stats:     *stats,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add follow relationship info if current user is authenticated
	if currentUserID != "" && currentUserID != userID {
		isFollowing, _ := s.checkIsFollowing(ctx, currentUserID, userID)
		isFollowedBy, _ := s.checkIsFollowing(ctx, userID, currentUserID)
		mutualCount, _ := s.getMutualConnectionsCount(ctx, currentUserID, userID)

		response.IsFollowing = isFollowing
		response.IsFollowedBy = isFollowedBy
		response.MutualConnectionsCount = mutualCount
	}

	return response, nil
}

// GetOwnProfile retrieves the authenticated user's full profile
func (s *profileService) GetOwnProfile(ctx context.Context, userID string) (*OwnProfileResponse, error) {
	// Get user from repository
	user, err := s.profileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get profile stats
	stats, err := s.getProfileStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile stats: %w", err)
	}

	// Build response with private fields
	response := &OwnProfileResponse{
		ID:                  user.ID,
		Username:            user.Username,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               user.Email,
		Phone:               user.Phone,
		DateOfBirth:         user.DOB.Format("2006-01-02"),
		AvatarURL:           user.AvatarURL,
		Bio:                 user.Bio,
		Stats:               *stats,
		EmailVerified:       user.EmailVerified,
		PhoneVerified:       user.PhoneVerified,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:           user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		PrivacySettings: PrivacySettings{
			WhoCanMessage:     user.WhoCanMessage,
			WhoCanSeePosts:    user.WhoCanSeePosts,
			ShowFollowersList: user.ShowFollowersList,
			ShowFollowingList: user.ShowFollowingList,
		},
	}

	return response, nil
}

// Helper methods

// getProfileStats gets user statistics (posts, followers, following)
func (s *profileService) getProfileStats(ctx context.Context, userID string) (*ProfileStats, error) {
	stats := &ProfileStats{}

	// Get posts count
	postsQuery := `SELECT COUNT(*) FROM posts WHERE user_id = $1 AND is_active = true`
	err := s.db.QueryRow(ctx, postsQuery, userID).Scan(&stats.PostsCount)
	if err != nil {
		stats.PostsCount = 0
	}

	// Get followers count
	followersQuery := `SELECT COUNT(*) FROM user_follows WHERE following_id = $1`
	err = s.db.QueryRow(ctx, followersQuery, userID).Scan(&stats.FollowersCount)
	if err != nil {
		stats.FollowersCount = 0
	}

	// Get following count
	followingQuery := `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`
	err = s.db.QueryRow(ctx, followingQuery, userID).Scan(&stats.FollowingCount)
	if err != nil {
		stats.FollowingCount = 0
	}

	return stats, nil
}

// checkIsFollowing checks if user1 follows user2
func (s *profileService) checkIsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_follows WHERE follower_id = $1 AND following_id = $2)`
	var exists bool
	err := s.db.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// getMutualConnectionsCount gets the count of mutual connections
func (s *profileService) getMutualConnectionsCount(ctx context.Context, user1ID, user2ID string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT uf1.following_id)
		FROM user_follows uf1
		INNER JOIN user_follows uf2 ON uf1.following_id = uf2.following_id
		WHERE uf1.follower_id = $1 AND uf2.follower_id = $2
		AND uf1.following_id NOT IN ($1, $2)
	`
	var count int
	err := s.db.QueryRow(ctx, query, user1ID, user2ID).Scan(&count)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

// Placeholder implementations for other methods

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error) {
	// TODO: Implement in Phase 4
	return nil, nil
}

func (s *profileService) UploadAvatar(ctx context.Context, userID string, fileBytes []byte, filename string) (*AvatarUploadResponse, error) {
	// TODO: Implement in Phase 5
	return nil, nil
}

func (s *profileService) DeleteAvatar(ctx context.Context, userID string) error {
	// TODO: Implement in Phase 5
	return nil
}

func (s *profileService) GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error) {
	// TODO: Implement in Phase 6
	return nil, nil
}

func (s *profileService) UpdatePrivacySettings(ctx context.Context, userID string, req *UpdatePrivacyRequest) (*PrivacySettings, error) {
	// TODO: Implement in Phase 6
	return nil, nil
}

func (s *profileService) GetMutualConnections(ctx context.Context, currentUserID, targetUserID string, page, limit int) (*MutualConnectionsResponse, error) {
	// TODO: Implement in Phase 7
	return nil, nil
}

func (s *profileService) GetMutualConnectionsCount(ctx context.Context, currentUserID, targetUserID string) (int, error) {
	// TODO: Implement in Phase 7
	return 0, nil
}
