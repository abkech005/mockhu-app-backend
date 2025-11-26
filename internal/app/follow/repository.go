package follow

import "context"

// FollowRepository defines the interface for follow data operations
type FollowRepository interface {
	// Follow creates a follow relationship
	Follow(ctx context.Context, followerID, followingID string) error

	// Unfollow removes a follow relationship
	Unfollow(ctx context.Context, followerID, followingID string) error

	// IsFollowing checks if follower follows following
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)

	// GetFollowers returns users who follow the given user
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*Follow, error)

	// GetFollowing returns users that the given user follows
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*Follow, error)

	// GetFollowerCount returns the number of followers
	GetFollowerCount(ctx context.Context, userID string) (int, error)

	// GetFollowingCount returns the number of users being followed
	GetFollowingCount(ctx context.Context, userID string) (int, error)

	// GetFollowStats returns both follower and following counts
	GetFollowStats(ctx context.Context, userID string) (*FollowStats, error)
}
