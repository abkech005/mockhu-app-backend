package messaging

import (
	"context"
	"fmt"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/follow"
)

// PrivacyChecker handles messaging privacy checks
type PrivacyChecker struct {
	userRepo   auth.UserRepository
	followRepo follow.FollowRepository
	blockRepo  BlockRepository
}

// NewPrivacyChecker creates a new privacy checker
func NewPrivacyChecker(userRepo auth.UserRepository, followRepo follow.FollowRepository, blockRepo BlockRepository) *PrivacyChecker {
	return &PrivacyChecker{
		userRepo:   userRepo,
		followRepo: followRepo,
		blockRepo:  blockRepo,
	}
}

// CanMessage checks if sender can message recipient
// Returns (canMessage bool, reason string, error)
func (p *PrivacyChecker) CanMessage(ctx context.Context, senderID, recipientID string, existingConversation bool) (bool, string, error) {
	// Cannot message yourself
	if senderID == recipientID {
		return false, "Cannot message yourself", nil
	}

	// Check if either user has blocked the other
	blocked, err := p.blockRepo.IsBlocked(ctx, senderID, recipientID)
	if err != nil {
		return false, "", fmt.Errorf("failed to check blocking status: %w", err)
	}
	if blocked {
		// Check which direction
		senderBlocked, err := p.blockRepo.IsUserBlocked(ctx, recipientID, senderID)
		if err != nil {
			return false, "", fmt.Errorf("failed to check if sender is blocked: %w", err)
		}
		if senderBlocked {
			return false, "You have been blocked by this user", nil
		}
		return false, "You have blocked this user", nil
	}

	// If there's an existing conversation, allow messaging (bypass privacy settings)
	if existingConversation {
		return true, "", nil
	}

	// Get recipient's privacy settings
	recipient, err := p.userRepo.FindByID(ctx, recipientID)
	if err != nil {
		return false, "", fmt.Errorf("failed to get recipient: %w", err)
	}

	// Check who_can_message setting
	switch recipient.WhoCanMessage {
	case "everyone":
		return true, "", nil

	case "followers":
		// Check if sender follows recipient
		isFollowing, err := p.followRepo.IsFollowing(ctx, senderID, recipientID)
		if err != nil {
			return false, "", fmt.Errorf("failed to check following status: %w", err)
		}
		if !isFollowing {
			return false, "Only followers can message this user", nil
		}
		return true, "", nil

	case "none":
		return false, "This user has disabled new messages", nil

	default:
		// Default to everyone if setting is invalid
		return true, "", nil
	}
}

// IsBlocked checks if either user has blocked the other
func (p *PrivacyChecker) IsBlocked(ctx context.Context, user1ID, user2ID string) (bool, error) {
	return p.blockRepo.IsBlocked(ctx, user1ID, user2ID)
}

