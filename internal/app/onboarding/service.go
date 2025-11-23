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
// Validates user, checks verification, updates profile, and marks onboarding complete
func (s *Service) CompleteOnboarding(ctx context.Context, req *CompleteOnboardingRequest) (*CompleteOnboardingResponse, error) {
	// 1. Get user by ID
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Check if email or phone is verified
	if !user.EmailVerified && !user.PhoneVerified {
		return nil, errors.New("please verify your email or phone before onboarding")
	}

	// 3. Check if username is already taken (if different from current)
	if req.Username != user.Username && req.Username != "" {
		existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("username already taken")
		}
	}

	// 4. Parse DOB (format: YYYY-MM-DD)
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// 5. Validate age (must be 13+)
	age := time.Now().Year() - dob.Year()
	if age < 13 {
		return nil, errors.New("you must be at least 13 years old")
	}
	if age > 120 {
		return nil, errors.New("invalid date of birth")
	}

	// 6. Update user profile
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = req.Username
	user.DOB = dob
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	user.UpdatedAt = time.Now()

	// 7. Mark onboarding as complete
	now := time.Now()
	user.OnboardingCompleted = true
	user.OnboardedAt = &now

	// 8. Save user to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// 9. Save interests (if provided)
	interestsCount := 0
	if len(req.Interests) > 0 {
		// Find interests by slugs
		interests, err := s.interestRepo.FindBySlugs(ctx, req.Interests)
		if err != nil {
			log.Printf("⚠️ Failed to find interests: %v", err)
			// Don't fail onboarding if interests fail
		} else if len(interests) > 0 {
			// Extract interest IDs
			interestIDs := make([]string, len(interests))
			for i, interest := range interests {
				interestIDs[i] = interest.ID
			}

			// Save user interests
			if err := s.interestRepo.AddUserInterests(ctx, user.ID, interestIDs); err != nil {
				log.Printf("⚠️ Failed to save user interests: %v", err)
				// Don't fail onboarding if interests fail
			} else {
				interestsCount = len(interests)
				log.Printf("✅ Saved %d interests for user %s", interestsCount, user.ID)
			}
		}
	}

	// 10. Return success response
	return &CompleteOnboardingResponse{
		Success:             true,
		Message:             "onboarding completed successfully",
		UserID:              user.ID,
		ProfileCompleted:    true,
		InterestsCount:      interestsCount,
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
