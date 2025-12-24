package postfeed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"mockhu-app-backend/internal/app/auth"
)

// Service handles postfeed business logic
type Service struct {
	repo     Repository
	authRepo auth.UserRepository
}

// NewService creates a new postfeed service
func NewService(repo Repository, authRepo auth.UserRepository) *Service {
	return &Service{repo: repo, authRepo: authRepo}
}

// Create creates a new postfeed
func (s *Service) Create(ctx context.Context, userID string, req CreatePostfeedRequest) (*PostfeedResponse, error) {
	// Validate type
	if !IsValidType(req.Type) {
		return nil, fmt.Errorf("invalid post type: %s", req.Type)
	}

	// Validate title
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 255 {
		return nil, errors.New("title must be 255 characters or less")
	}

	// Validate and set visibility (default to public)
	visibility := req.Visibility
	if visibility == "" {
		visibility = VisibilityPublic
	}
	if visibility != VisibilityPublic && visibility != VisibilityPrivate && visibility != VisibilityFollowersOnly {
		return nil, fmt.Errorf("invalid visibility: %s (must be public, private, or followers_only)", visibility)
	}

	// Validate type-specific metadata
	if err := s.validateMetadata(req.Type, req.Metadata); err != nil {
		return nil, err
	}

	postfeed := &Postfeed{
		UserID:      userID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		Visibility:  visibility,
		IsAnonymous: req.IsAnonymous,
		Metadata:    req.Metadata,
	}

	if err := s.repo.Create(ctx, postfeed); err != nil {
		return nil, fmt.Errorf("failed to create postfeed: %w", err)
	}

	return s.toResponse(ctx, postfeed)
}

// GetByID retrieves a postfeed by ID and increments view count
func (s *Service) GetByID(ctx context.Context, id string) (*PostfeedResponse, error) {
	postfeed, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously (fire and forget)
	go func() {
		_ = s.repo.IncrementViewCount(context.Background(), id)
	}()

	return s.toResponse(ctx, postfeed)
}

// Update updates a postfeed
func (s *Service) Update(ctx context.Context, id string, userID string, req UpdatePostfeedRequest) (*PostfeedResponse, error) {
	// Get existing postfeed
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if existing.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	updates := make(map[string]interface{})

	if req.Title != "" {
		if len(req.Title) > 255 {
			return nil, errors.New("title must be 255 characters or less")
		}
		updates["title"] = req.Title
	}

	if req.Content != "" {
		updates["content"] = req.Content
	}

	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	if req.Visibility != "" {
		if req.Visibility != VisibilityPublic && req.Visibility != VisibilityPrivate && req.Visibility != VisibilityFollowersOnly {
			return nil, fmt.Errorf("invalid visibility: %s", req.Visibility)
		}
		updates["visibility"] = req.Visibility
	}

	if req.Metadata != nil {
		if err := s.validateMetadata(existing.Type, req.Metadata); err != nil {
			return nil, err
		}
		updates["metadata"] = req.Metadata
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	// Fetch updated postfeed
	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(ctx, updated)
}

// Delete soft deletes a postfeed
func (s *Service) Delete(ctx context.Context, id string, userID string) error {
	// Get existing postfeed
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if existing.UserID != userID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	return s.repo.Delete(ctx, id)
}

// List retrieves postfeeds with filters
func (s *Service) List(ctx context.Context, req ListPostfeedsRequest) (*ListPostfeedsResponse, error) {
	// Defaults
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.Limit

	filter := ListFilter{
		Type:   req.Type,
		UserID: req.UserID,
		Tag:    req.Tag,
		Limit:  req.Limit,
		Offset: offset,
	}

	postfeeds, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]PostfeedResponse, len(postfeeds))
	for i, p := range postfeeds {
		resp, err := s.toResponse(ctx, &p)
		if err != nil {
			continue
		}
		responses[i] = *resp
	}

	totalPages := (total + req.Limit - 1) / req.Limit

	return &ListPostfeedsResponse{
		Postfeeds:  responses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserPostfeeds retrieves postfeeds by user
func (s *Service) GetUserPostfeeds(ctx context.Context, userID string, page, limit int) (*ListPostfeedsResponse, error) {
	return s.List(ctx, ListPostfeedsRequest{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
}

// validateMetadata validates type-specific metadata
func (s *Service) validateMetadata(postType string, metadata json.RawMessage) error {
	if len(metadata) == 0 {
		return nil
	}

	switch postType {
	case TypeDoubt:
		var m DoubtMetadata
		if err := json.Unmarshal(metadata, &m); err != nil {
			return fmt.Errorf("invalid doubt metadata: %w", err)
		}
	case TypeQuiz:
		var m QuizMetadata
		if err := json.Unmarshal(metadata, &m); err != nil {
			return fmt.Errorf("invalid quiz metadata: %w", err)
		}
		if len(m.Questions) == 0 {
			return errors.New("quiz must have at least one question")
		}
		for i, q := range m.Questions {
			if q.Question == "" {
				return fmt.Errorf("question %d is empty", i+1)
			}
			if len(q.Options) < 2 {
				return fmt.Errorf("question %d must have at least 2 options", i+1)
			}
			if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Options) {
				return fmt.Errorf("question %d has invalid correct_index", i+1)
			}
		}
	case TypeProgress:
		var m ProgressMetadata
		if err := json.Unmarshal(metadata, &m); err != nil {
			return fmt.Errorf("invalid progress metadata: %w", err)
		}
	case TypeResource:
		var m ResourceMetadata
		if err := json.Unmarshal(metadata, &m); err != nil {
			return fmt.Errorf("invalid resource metadata: %w", err)
		}
	}

	return nil
}

// toResponse converts Postfeed to PostfeedResponse with author info
func (s *Service) toResponse(ctx context.Context, p *Postfeed) (*PostfeedResponse, error) {
	resp := &PostfeedResponse{
		ID:           p.ID,
		UserID:       p.UserID,
		Type:         p.Type,
		Title:        p.Title,
		Content:      p.Content,
		Tags:         p.Tags,
		Visibility:   p.Visibility,
		IsAnonymous:  p.IsAnonymous,
		Metadata:     p.Metadata,
		ViewCount:    p.ViewCount,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ShareCount:   p.ShareCount,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add author info if not anonymous
	if !p.IsAnonymous && s.authRepo != nil {
		user, err := s.authRepo.FindByID(ctx, p.UserID)
		if err == nil && user != nil {
			resp.Author = &AuthorInfo{
				ID:        user.ID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				AvatarURL: user.AvatarURL,
			}
		}
	}

	return resp, nil
}
