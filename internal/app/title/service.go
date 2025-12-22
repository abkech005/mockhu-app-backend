package title

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service handles title business logic
type Service struct {
	repo TitleRepository
}

// NewService creates a new title service
func NewService(repo TitleRepository) *Service {
	return &Service{repo: repo}
}

// GetAllTitles retrieves all available titles
func (s *Service) GetAllTitles(ctx context.Context) ([]Title, error) {
	titles, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get titles: %w", err)
	}

	if titles == nil {
		titles = []Title{}
	}

	return titles, nil
}

// GetTitlesByDefinedBy retrieves titles filtered by defined_by (admin or user)
func (s *Service) GetTitlesByDefinedBy(ctx context.Context, definedBy string) ([]Title, error) {
	// Validate defined_by
	if definedBy != DefinedByAdmin && definedBy != DefinedByUser {
		return nil, fmt.Errorf("invalid defined_by value: %s (must be 'admin' or 'user')", definedBy)
	}

	return s.repo.FindByDefinedBy(ctx, definedBy)
}

// GetTitleByID retrieves a title by its ID
func (s *Service) GetTitleByID(ctx context.Context, id string) (*Title, error) {
	if id == "" {
		return nil, errors.New("title ID is required")
	}

	return s.repo.FindByID(ctx, id)
}

// GetTitleByName retrieves a title by its name
func (s *Service) GetTitleByName(ctx context.Context, name string) (*Title, error) {
	if name == "" {
		return nil, errors.New("title name is required")
	}

	return s.repo.FindByName(ctx, name)
}

// CreateTitle creates a new title (user-defined)
func (s *Service) CreateTitle(ctx context.Context, name, description string) (*Title, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("title name is required")
	}

	if len(name) > 100 {
		return nil, errors.New("title name must be 100 characters or less")
	}

	// Check if title with same name already exists
	existing, _ := s.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, fmt.Errorf("title with name '%s' already exists", name)
	}

	title := &Title{
		Name:        name,
		Description: strings.TrimSpace(description),
		DefinedBy:   DefinedByUser, // User-created titles
		UsedByCount: 0,
	}

	if err := s.repo.Create(ctx, title); err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}

	return title, nil
}

// CreateAdminTitle creates a new admin-defined title
func (s *Service) CreateAdminTitle(ctx context.Context, name, description string) (*Title, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("title name is required")
	}

	if len(name) > 100 {
		return nil, errors.New("title name must be 100 characters or less")
	}

	// Check if title with same name already exists
	existing, _ := s.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, fmt.Errorf("title with name '%s' already exists", name)
	}

	title := &Title{
		Name:        name,
		Description: strings.TrimSpace(description),
		DefinedBy:   DefinedByAdmin, // Admin-created titles
		UsedByCount: 0,
	}

	if err := s.repo.Create(ctx, title); err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}

	return title, nil
}

// UpdateTitle updates an existing title
func (s *Service) UpdateTitle(ctx context.Context, id, name, description string) (*Title, error) {
	if id == "" {
		return nil, errors.New("title ID is required")
	}

	// Get existing title
	title, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		name = strings.TrimSpace(name)
		if len(name) > 100 {
			return nil, errors.New("title name must be 100 characters or less")
		}
		// Check if another title with same name exists
		existing, _ := s.repo.FindByName(ctx, name)
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("title with name '%s' already exists", name)
		}
		title.Name = name
	}

	if description != "" {
		title.Description = strings.TrimSpace(description)
	}

	if err := s.repo.Update(ctx, title); err != nil {
		return nil, fmt.Errorf("failed to update title: %w", err)
	}

	return title, nil
}

// DeleteTitle deletes a title
func (s *Service) DeleteTitle(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("title ID is required")
	}

	// Check if title exists and is not being used
	title, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if title.UsedByCount > 0 {
		return fmt.Errorf("cannot delete title '%s' as it is being used by %d users", title.Name, title.UsedByCount)
	}

	return s.repo.Delete(ctx, id)
}

// IncrementUsage increments the used_by_count for a title
func (s *Service) IncrementUsage(ctx context.Context, id string) (int, error) {
	if id == "" {
		return 0, errors.New("title ID is required")
	}

	return s.repo.IncrementUsedByCount(ctx, id)
}

// DecrementUsage decrements the used_by_count for a title
func (s *Service) DecrementUsage(ctx context.Context, id string) (int, error) {
	if id == "" {
		return 0, errors.New("title ID is required")
	}

	return s.repo.DecrementUsedByCount(ctx, id)
}

// GetMostUsedTitles returns the most popular titles
func (s *Service) GetMostUsedTitles(ctx context.Context, limit int) ([]Title, error) {
	if limit <= 0 {
		limit = 10
	}

	titles, err := s.repo.GetMostUsedTitles(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used titles: %w", err)
	}

	if titles == nil {
		titles = []Title{}
	}

	return titles, nil
}
