package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"mockhu-app-backend/internal/app/interest"
	"mockhu-app-backend/internal/app/location"
	"mockhu-app-backend/internal/app/title"
	"mockhu-app-backend/internal/pkg/avatar"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/microcosm-cc/bluemonday"
)

// ProfileService defines the business logic for profile operations
type ProfileService interface {
	// Profile viewing
	GetUserProfile(ctx context.Context, userID, currentUserID string) (*ProfileResponse, error)
	GetOwnProfile(ctx context.Context, userID string) (*OwnProfileResponse, error)

	// Profile management
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error)
	UpdateFullProfile(ctx context.Context, userID string, req *UpdateFullProfileRequest) (*ProfileResponse, error)

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
	profileRepo     ProfileRepository
	titleRepo       title.TitleRepository
	locationRepo    location.LocationRepository
	interestService *interest.Service // Using concrete struct as dependency
	db              *pgxpool.Pool
	policy          *bluemonday.Policy
}

// NewService creates a new profile service
func NewService(
	profileRepo ProfileRepository,
	titleRepo title.TitleRepository,
	locationRepo location.LocationRepository,
	interestService *interest.Service,
	db *pgxpool.Pool,
) ProfileService {
	return &profileService{
		profileRepo:     profileRepo,
		titleRepo:       titleRepo,
		locationRepo:    locationRepo,
		interestService: interestService,
		db:              db,
		policy:          bluemonday.UGCPolicy(),
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

	// Get user interests
	interests, _ := s.interestService.GetUserInterests(ctx, userID)

	// Build response
	response := &ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Title:     user.Title,
		Place:     user.Place,
		Interests: interests,
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

	// Get user interests
	interests, _ := s.interestService.GetUserInterests(ctx, userID)

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
		Title:               user.Title,
		Place:               user.Place,
		Interests:           interests,
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
	// Validate input
	if err := s.validateUpdateProfileRequest(req); err != nil {
		return nil, err
	}

	// Fetch current profile to handle count decrements
	currentUser, err := s.profileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current profile: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Bio != "" {
		// Sanitize bio
		sanitizedBio := s.policy.Sanitize(req.Bio)
		updates["bio"] = sanitizedBio
	}
	if req.Username != "" && req.Username != currentUser.Username {
		// Check username uniqueness
		exists, err := s.profileRepo.CheckUsernameExists(ctx, req.Username, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if exists {
			return nil, errors.New("username already taken")
		}
		updates["username"] = req.Username
	}

	// Handle Title logic
	if req.Title != "" {
		// Update string column
		sanitizedTitle := s.policy.Sanitize(req.Title)
		updates["title"] = sanitizedTitle

		// Find or create title
		titleObj, err := s.resolveTitle(ctx, sanitizedTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve title: %w", err)
		}
		if titleObj != nil {
			updates["title_id"] = titleObj.ID

			// Decrement old title count if changed
			if currentUser.TitleID != nil && *currentUser.TitleID != titleObj.ID {
				_, _ = s.titleRepo.DecrementUsedByCount(ctx, *currentUser.TitleID)
			}
			// Increment new title count
			if currentUser.TitleID == nil || *currentUser.TitleID != titleObj.ID {
				_, _ = s.titleRepo.IncrementUsedByCount(ctx, titleObj.ID)
			}
		}
	}

	// Handle Place logic
	if req.Place != "" {
		// Update string column
		updates["place"] = req.Place

		// Find or create location
		locObj, err := s.resolveLocation(ctx, req.Place)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve location: %w", err)
		}

		if locObj != nil {
			updates["location_id"] = locObj.ID

			// Decrement old location count if changed
			if currentUser.LocationID != nil && *currentUser.LocationID != locObj.ID {
				_, _ = s.locationRepo.DecrementUsedByCount(ctx, *currentUser.LocationID)
			}
			// Increment new location count
			if currentUser.LocationID == nil || *currentUser.LocationID != locObj.ID {
				_, _ = s.locationRepo.IncrementUsedByCount(ctx, locObj.ID)
			}
		} else {
			// Location not resolved, clear the link
			updates["location_id"] = nil

			// Decrement old location count if it existed
			if currentUser.LocationID != nil {
				_, _ = s.locationRepo.DecrementUsedByCount(ctx, *currentUser.LocationID)
			}
		}
	}

	// If no updates, return (or return current profile)
	if len(updates) == 0 {
		return s.GetUserProfile(ctx, userID, userID)
	}

	// Update profile
	err = s.profileRepo.UpdateProfile(ctx, userID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Get updated profile
	return s.GetUserProfile(ctx, userID, userID)
}

// UpdateFullProfile handles composite profile and interests update
func (s *profileService) UpdateFullProfile(ctx context.Context, userID string, req *UpdateFullProfileRequest) (*ProfileResponse, error) {
	// 1. Update Profile Fields
	profileResp, err := s.UpdateProfile(ctx, userID, &req.UpdateProfileRequest)
	if err != nil {
		return nil, err
	}

	// 2. Update Interests
	if req.InterestSlugs != nil {
		_, err := s.interestService.ReplaceUserInterests(ctx, userID, req.InterestSlugs)
		if err != nil {
			return nil, fmt.Errorf("failed to update interests: %w", err)
		}
	}

	return profileResp, nil
}

// resolveTitle finds or creates a title
func (s *profileService) resolveTitle(ctx context.Context, name string) (*title.Title, error) {
	existing, err := s.titleRepo.FindByName(ctx, name)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new title
	newTitle := &title.Title{
		Name:      name,
		DefinedBy: "user",
	}
	if err := s.titleRepo.Create(ctx, newTitle); err != nil {
		return nil, err
	}
	return newTitle, nil
}

// resolveLocation finds or creates a location
func (s *profileService) resolveLocation(ctx context.Context, place string) (*location.Location, error) {
	// Simple comma splitting: "City, Country"
	parts := strings.Split(place, ",")
	if len(parts) < 2 {
		// Try search by entire string
		results, err := s.locationRepo.Search(ctx, place)
		if err == nil && len(results) > 0 {
			return &results[0], nil
		}
		return nil, nil // Could not resolve
	}

	city := strings.TrimSpace(parts[0])
	country := strings.TrimSpace(parts[len(parts)-1]) // Last part is country

	// Try find exact
	existing, err := s.locationRepo.FindByCityAndCountry(ctx, city, country)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new location
	newLoc := &location.Location{
		City:    city,
		Country: country,
	}
	if err := s.locationRepo.Create(ctx, newLoc); err != nil {
		return nil, err
	}
	return newLoc, nil
}

// validateUpdateProfileRequest validates the update profile request
func (s *profileService) validateUpdateProfileRequest(req *UpdateProfileRequest) error {
	if req.FirstName != "" {
		if len(req.FirstName) < 1 || len(req.FirstName) > 50 {
			return errors.New("first name must be between 1 and 50 characters")
		}
	}

	if req.LastName != "" {
		if len(req.LastName) < 1 || len(req.LastName) > 50 {
			return errors.New("last name must be between 1 and 50 characters")
		}
	}

	if req.Bio != "" {
		if len(req.Bio) > 500 {
			return errors.New("bio must not exceed 500 characters")
		}
	}

	if req.Username != "" {
		if len(req.Username) < 3 || len(req.Username) > 30 {
			return errors.New("username must be between 3 and 30 characters")
		}
		// Check if username contains only alphanumeric and underscore
		for _, char := range req.Username {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '_') {
				return errors.New("username can only contain letters, numbers, and underscores")
			}
		}
	}

	return nil
}

func (s *profileService) UploadAvatar(ctx context.Context, userID string, fileBytes []byte, filename string) (*AvatarUploadResponse, error) {
	// Get current user to retrieve old avatar
	user, err := s.profileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Process and save new avatar (resize, validate, save to storage)
	avatarURL, err := avatar.ProcessAndSave(fileBytes, filename)
	if err != nil {
		// Return user-friendly error messages
		if errors.Is(err, avatar.ErrFileTooBig) {
			return nil, errors.New("file size exceeds 5MB")
		}
		if errors.Is(err, avatar.ErrInvalidFileType) {
			return nil, errors.New("invalid file type, only JPEG, PNG, and WebP allowed")
		}
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	// Update database with new avatar URL
	err = s.profileRepo.UpdateAvatar(ctx, userID, avatarURL)
	if err != nil {
		// If database update fails, try to delete the uploaded file
		_ = avatar.DeleteAvatar(avatarURL)
		return nil, fmt.Errorf("failed to update avatar: %w", err)
	}

	// Delete old avatar file if it exists
	if user.AvatarURL != "" {
		_ = avatar.DeleteAvatar(user.AvatarURL)
	}

	return &AvatarUploadResponse{
		AvatarURL: avatarURL,
		Message:   "avatar uploaded successfully",
	}, nil
}

func (s *profileService) DeleteAvatar(ctx context.Context, userID string) error {
	// Get current user to retrieve avatar URL
	user, err := s.profileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete avatar file
	if user.AvatarURL != "" {
		err = avatar.DeleteAvatar(user.AvatarURL)
		if err != nil {
			return fmt.Errorf("failed to delete avatar file: %w", err)
		}
	}

	// Update database (set avatar_url to empty string)
	err = s.profileRepo.UpdateAvatar(ctx, userID, "")
	if err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}

func (s *profileService) GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error) {
	// Get privacy settings from repository
	settings, err := s.profileRepo.GetPrivacySettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get privacy settings: %w", err)
	}

	return settings, nil
}

func (s *profileService) UpdatePrivacySettings(ctx context.Context, userID string, req *UpdatePrivacyRequest) (*PrivacySettings, error) {
	// Validate privacy settings
	if err := s.validatePrivacySettings(req); err != nil {
		return nil, err
	}

	// Update privacy settings
	err := s.profileRepo.UpdatePrivacySettings(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update privacy settings: %w", err)
	}

	// Get and return updated settings
	return s.GetPrivacySettings(ctx, userID)
}

// validatePrivacySettings validates privacy settings request
func (s *profileService) validatePrivacySettings(req *UpdatePrivacyRequest) error {
	// Validate who_can_message
	if req.WhoCanMessage != "" {
		if req.WhoCanMessage != "everyone" && req.WhoCanMessage != "followers" && req.WhoCanMessage != "none" {
			return errors.New("who_can_message must be 'everyone', 'followers', or 'none'")
		}
	}

	// Validate who_can_see_posts
	if req.WhoCanSeePosts != "" {
		if req.WhoCanSeePosts != "everyone" && req.WhoCanSeePosts != "followers" && req.WhoCanSeePosts != "none" {
			return errors.New("who_can_see_posts must be 'everyone', 'followers', or 'none'")
		}
	}

	return nil
}

func (s *profileService) GetMutualConnections(ctx context.Context, currentUserID, targetUserID string, page, limit int) (*MutualConnectionsResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get mutual connections from repository
	users, err := s.profileRepo.GetMutualConnections(ctx, currentUserID, targetUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get mutual connections: %w", err)
	}

	// Convert to response format
	connections := make([]*MutualConnectionUser, 0, len(users))
	for _, user := range users {
		connections = append(connections, &MutualConnectionUser{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			AvatarURL: user.AvatarURL,
		})
	}

	// Get total count for pagination
	totalCount, err := s.profileRepo.GetMutualConnectionsCount(ctx, currentUserID, targetUserID)
	if err != nil {
		totalCount = len(connections) // Fallback to current page count
	}

	totalPages := (totalCount + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &MutualConnectionsResponse{
		MutualConnections: connections,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: totalCount,
			Limit:      limit,
		},
	}, nil
}

func (s *profileService) GetMutualConnectionsCount(ctx context.Context, currentUserID, targetUserID string) (int, error) {
	count, err := s.profileRepo.GetMutualConnectionsCount(ctx, currentUserID, targetUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to get mutual connections count: %w", err)
	}

	return count, nil
}
