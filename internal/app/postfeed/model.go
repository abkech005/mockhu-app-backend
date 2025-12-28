package postfeed

import (
	"encoding/json"
	"time"
)

// Post type constants
const (
	TypeDoubt    = "doubt"
	TypeQuiz     = "quiz"
	TypeProgress = "progress"
	TypeResource = "resource"
)

// Visibility constants
const (
	VisibilityPublic        = "public"
	VisibilityPrivate       = "private"
	VisibilityFollowersOnly = "followers"
)

// Postfeed represents a social media feed post
type Postfeed struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Content      string          `json:"content,omitempty"`
	Tags         []string        `json:"tags,omitempty"`
	Media        []MediaItem     `json:"media,omitempty"`
	Visibility   string          `json:"visibility"`
	IsAnonymous  bool            `json:"is_anonymous"`
	IsActive     bool            `json:"is_active"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	ViewCount    int             `json:"view_count"`
	LikeCount    int             `json:"like_count"`
	CommentCount int             `json:"comment_count"`
	ShareCount   int             `json:"share_count"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// MediaItem represents a media attachment (image, video, etc.)
type MediaItem struct {
	URL          string `json:"url"`
	Type         string `json:"type"` // image, video, audio, document
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	Duration     int    `json:"duration,omitempty"` // for video/audio in seconds
	FileName     string `json:"file_name,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// DoubtMetadata for doubt-type posts
type DoubtMetadata struct {
	Subject      string  `json:"subject,omitempty"`
	IsSolved     bool    `json:"is_solved"`
	BestAnswerID *string `json:"best_answer_id,omitempty"`
}

// QuizQuestion represents a single quiz question
type QuizQuestion struct {
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correct_index"`
}

// QuizMetadata for quiz-type posts
type QuizMetadata struct {
	Questions        []QuizQuestion `json:"questions"`
	TimeLimitSeconds int            `json:"time_limit_seconds,omitempty"`
	Difficulty       string         `json:"difficulty,omitempty"` // easy, medium, hard
}

// ProgressMetadata for progress-type posts
type ProgressMetadata struct {
	Milestone   string `json:"milestone,omitempty"`
	Percentage  int    `json:"percentage,omitempty"`
	StreakDays  int    `json:"streak_days,omitempty"`
	BadgeEarned string `json:"badge_earned,omitempty"`
}

// ResourceMetadata for resource-type posts
type ResourceMetadata struct {
	ResourceType    string `json:"resource_type,omitempty"` // video, pdf, link, article
	URL             string `json:"url,omitempty"`
	FileURL         string `json:"file_url,omitempty"`
	Platform        string `json:"platform,omitempty"` // youtube, notion, etc.
	DurationMinutes int    `json:"duration_minutes,omitempty"`
}

// ValidTypes returns all valid post types
func ValidTypes() []string {
	return []string{TypeDoubt, TypeQuiz, TypeProgress, TypeResource}
}

// IsValidType checks if a type is valid
func IsValidType(t string) bool {
	for _, valid := range ValidTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// Like represents a like on a postfeed
type Like struct {
	ID         string    `json:"id"`
	PostfeedID string    `json:"postfeed_id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Comment represents a comment on a postfeed
type Comment struct {
	ID         string    `json:"id"`
	PostfeedID string    `json:"postfeed_id"`
	UserID     string    `json:"user_id"`
	ParentID   *string   `json:"parent_id,omitempty"`
	Content    string    `json:"content"`
	IsActive   bool      `json:"is_active"`
	LikeCount  int       `json:"like_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Share represents a share of a postfeed
type Share struct {
	ID         string    `json:"id"`
	PostfeedID string    `json:"postfeed_id"`
	UserID     string    `json:"user_id"`
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
