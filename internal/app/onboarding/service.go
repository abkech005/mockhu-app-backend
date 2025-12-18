package onboarding

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/interest"
)

// Service handles onboarding business logic
type Service struct {
	userRepo     auth.UserRepository
	interestRepo interest.InterestRepository
}

// NewService creates a new onboarding service
func NewService(userRepo auth.UserRepository, interestRepo interest.InterestRepository) *Service {
	return &Service{
		userRepo:     userRepo,
		interestRepo: interestRepo,
	}
}

// CompleteOnboarding handles the entire onboarding process
// Validates user, updates profile, and marks onboarding complete
func (s *Service) CompleteOnboarding(ctx context.Context, req *CompleteOnboardingRequest) (*CompleteOnboardingResponse, error) {
	// 1. Get user by ID
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Check if username is already taken (if different from current)
	if req.Username != user.Username && req.Username != "" {
		existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("username already taken")
		}
	}

	// 4. Update user profile
	user.Username = req.Username
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	user.Bio = req.Bio
	user.Place = req.Place
	user.UpdatedAt = time.Now()

	// 5. Mark onboarding as complete
	now := time.Now()
	user.OnboardingCompleted = true
	user.OnboardedAt = &now

	// 6. Save user to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	log.Printf("✅ Onboarding completed for user %s", user.ID)

	// 7. Return success response
	return &CompleteOnboardingResponse{
		Success:             true,
		Message:             "onboarding completed successfully",
		UserID:              user.ID,
		OnboardingCompleted: true,
		OnboardedAt:         now,
	}, nil
}

// GetOnboardingStatus returns the current onboarding status for a user
func (s *Service) GetOnboardingStatus(ctx context.Context, userID string) (*OnboardingStatusResponse, error) {
	// Get user from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if profile is completed
	profileCompleted := user.FirstName != "" && user.LastName != "" && user.Username != ""

	// Build response
	response := &OnboardingStatusResponse{
		UserID:              user.ID,
		Email:               user.Email,
		EmailVerified:       user.EmailVerified,
		PhoneVerified:       user.PhoneVerified,
		ProfileCompleted:    profileCompleted,
		OnboardingCompleted: user.OnboardingCompleted,
	}

	// Determine next step
	if !user.EmailVerified && !user.PhoneVerified {
		response.NextStep = "verify_email_or_phone"
	} else if !profileCompleted {
		response.NextStep = "complete_profile"
	} else if !user.OnboardingCompleted {
		response.NextStep = "finalize_onboarding"
	} else {
		response.NextStep = "completed"
	}

	return response, nil
}
