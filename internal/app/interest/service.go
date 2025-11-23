package interest

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service handles interest business logic
type Service struct {
	repo InterestRepository
}

// NewService creates a new interest service
func NewService(repo InterestRepository) *Service {
	return &Service{repo: repo}
}

// GetAllInterests retrieves all available interests
func (s *Service) GetAllInterests(ctx context.Context) ([]Interest, error) {
	interests, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get interests: %w", err)
	}
	
	if interests == nil {
		interests = []Interest{}
	}
	
	return interests, nil
}

// GetInterestsByCategory retrieves interests filtered by category
func (s *Service) GetInterestsByCategory(ctx context.Context, category string) ([]Interest, error) {
	// Validate category
	validCategories := []string{
		CategoryTechnology, CategoryArts, CategorySports,
		CategoryEntertainment, CategoryLifestyle, CategoryBusiness,
		CategoryEducation, CategorySocial,
	}
	
	isValid := false
	for _, valid := range validCategories {
		if category == valid {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return nil, fmt.Errorf("invalid category: %s", category)
	}
	
	return s.repo.FindByCategory(ctx, category)
}

// GetCategories returns all categories with their interest counts
func (s *Service) GetCategories(ctx context.Context) ([]CategoryInfo, error) {
	counts, err := s.repo.CountByCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	
	categories := []CategoryInfo{}
	for name, count := range counts {
		categories = append(categories, CategoryInfo{
			Name:  name,
			Count: count,
		})
	}
	
	return categories, nil
}

// GetUserInterests retrieves all interests for a specific user
func (s *Service) GetUserInterests(ctx context.Context, userID string) ([]Interest, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	interests, err := s.repo.GetUserInterests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user interests: %w", err)
	}
	
	if interests == nil {
		interests = []Interest{}
	}
	
	return interests, nil
}

// AddUserInterests adds interests to a user by slugs
func (s *Service) AddUserInterests(ctx context.Context, userID string, slugs []string) ([]Interest, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	if len(slugs) == 0 {
		return nil, errors.New("at least one interest slug is required")
	}
	
	// Normalize slugs (lowercase, trim)
	normalizedSlugs := make([]string, len(slugs))
	for i, slug := range slugs {
		normalizedSlugs[i] = strings.ToLower(strings.TrimSpace(slug))
	}
	
	// Find interests by slugs
	interests, err := s.repo.FindBySlugs(ctx, normalizedSlugs)
	if err != nil {
		return nil, fmt.Errorf("failed to find interests: %w", err)
	}
	
	if len(interests) == 0 {
		return nil, errors.New("no valid interests found")
	}
	
	// Extract interest IDs
	interestIDs := make([]string, len(interests))
	for i, interest := range interests {
		interestIDs[i] = interest.ID
	}
	
	// Add interests to user
	if err := s.repo.AddUserInterests(ctx, userID, interestIDs); err != nil {
		return nil, fmt.Errorf("failed to add user interests: %w", err)
	}
	
	return interests, nil
}

// RemoveUserInterest removes an interest from a user
func (s *Service) RemoveUserInterest(ctx context.Context, userID string, slug string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	
	if slug == "" {
		return errors.New("interest slug is required")
	}
	
	// Find interest by slug
	interest, err := s.repo.FindBySlug(ctx, strings.ToLower(strings.TrimSpace(slug)))
	if err != nil {
		return fmt.Errorf("interest not found: %w", err)
	}
	
	// Remove interest from user
	if err := s.repo.RemoveUserInterest(ctx, userID, interest.ID); err != nil {
		return fmt.Errorf("failed to remove user interest: %w", err)
	}
	
	return nil
}

// ReplaceUserInterests replaces all user interests with new ones
func (s *Service) ReplaceUserInterests(ctx context.Context, userID string, slugs []string) ([]Interest, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	// Allow empty slugs to clear all interests
	if len(slugs) == 0 {
		if err := s.repo.ReplaceUserInterests(ctx, userID, []string{}); err != nil {
			return nil, fmt.Errorf("failed to clear user interests: %w", err)
		}
		return []Interest{}, nil
	}
	
	// Normalize slugs
	normalizedSlugs := make([]string, len(slugs))
	for i, slug := range slugs {
		normalizedSlugs[i] = strings.ToLower(strings.TrimSpace(slug))
	}
	
	// Find interests by slugs
	interests, err := s.repo.FindBySlugs(ctx, normalizedSlugs)
	if err != nil {
		return nil, fmt.Errorf("failed to find interests: %w", err)
	}
	
	if len(interests) == 0 {
		return nil, errors.New("no valid interests found")
	}
	
	// Extract interest IDs
	interestIDs := make([]string, len(interests))
	for i, interest := range interests {
		interestIDs[i] = interest.ID
	}
	
	// Replace user interests
	if err := s.repo.ReplaceUserInterests(ctx, userID, interestIDs); err != nil {
		return nil, fmt.Errorf("failed to replace user interests: %w", err)
	}
	
	return interests, nil
}

// CreateInterest creates a new interest (admin function)
func (s *Service) CreateInterest(ctx context.Context, name, slug, category, icon string) (*Interest, error) {
	// Validate input
	if name == "" || slug == "" || category == "" {
		return nil, errors.New("name, slug, and category are required")
	}
	
	interest := &Interest{
		Name:     name,
		Slug:     strings.ToLower(strings.TrimSpace(slug)),
		Category: strings.ToLower(strings.TrimSpace(category)),
		Icon:     icon,
	}
	
	if err := s.repo.Create(ctx, interest); err != nil {
		return nil, fmt.Errorf("failed to create interest: %w", err)
	}
	
	return interest, nil
}

