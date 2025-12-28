package suggestion

import (
	"context"
)

// Service handles suggestion business logic
type Service struct {
	repo *Repository
}

// NewService creates a new suggestion service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetUserSuggestions returns suggested users based on shared interests
func (s *Service) GetUserSuggestions(ctx context.Context, userID string, limit int) (*GetUserSuggestionsResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10 // Default limit
	}

	suggestions, err := s.repo.GetSuggestedUsers(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	if suggestions == nil {
		suggestions = []SuggestedUser{}
	}

	return &GetUserSuggestionsResponse{
		Suggestions: suggestions,
		Total:       len(suggestions),
	}, nil
}
