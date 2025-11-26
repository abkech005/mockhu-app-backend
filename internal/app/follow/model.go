package follow

import "time"

// Follow represents a follow relationship between two users
type Follow struct {
	ID          string    `json:"id"`
	FollowerID  string    `json:"follower_id"`  // User who follows
	FollowingID string    `json:"following_id"` // User being followed
	CreatedAt   time.Time `json:"created_at"`
}

// FollowStats contains follower/following counts for a user
type FollowStats struct {
	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`
}
