package location

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service handles location business logic
type Service struct {
	repo LocationRepository
}

// NewService creates a new location service
func NewService(repo LocationRepository) *Service {
	return &Service{repo: repo}
}

// GetAllLocations retrieves all locations
func (s *Service) GetAllLocations(ctx context.Context) ([]Location, error) {
	locations, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	if locations == nil {
		locations = []Location{}
	}

	return locations, nil
}

// GetLocationByID retrieves a location by ID
func (s *Service) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	if id == "" {
		return nil, errors.New("location ID is required")
	}

	return s.repo.FindByID(ctx, id)
}

// SearchLocations searches locations by query (autocomplete)
func (s *Service) SearchLocations(ctx context.Context, query string) ([]Location, error) {
	if query == "" {
		return s.repo.FindAll(ctx)
	}

	locations, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search locations: %w", err)
	}

	if locations == nil {
		locations = []Location{}
	}

	return locations, nil
}

// GetLocationsByCountry retrieves locations by country
func (s *Service) GetLocationsByCountry(ctx context.Context, country string) ([]Location, error) {
	if country == "" {
		return nil, errors.New("country is required")
	}

	return s.repo.FindByCountry(ctx, country)
}

// GetMostUsedLocations returns popular locations
func (s *Service) GetMostUsedLocations(ctx context.Context, limit int) ([]Location, error) {
	if limit <= 0 {
		limit = 10
	}

	locations, err := s.repo.GetMostUsedLocations(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular locations: %w", err)
	}

	if locations == nil {
		locations = []Location{}
	}

	return locations, nil
}

// CreateLocation creates a new location
func (s *Service) CreateLocation(ctx context.Context, city, country string) (*Location, error) {
	city = strings.TrimSpace(city)
	country = strings.TrimSpace(country)

	if city == "" {
		return nil, errors.New("city is required")
	}

	if country == "" {
		return nil, errors.New("country is required")
	}

	// Check if location already exists
	existing, err := s.repo.FindByCityAndCountry(ctx, city, country)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing location: %w", err)
	}

	if existing != nil {
		return existing, nil // Return existing location
	}

	location := &Location{
		City:        city,
		Country:     country,
		UsedByCount: 0,
	}

	if err := s.repo.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

// UpdateLocation updates an existing location
func (s *Service) UpdateLocation(ctx context.Context, id string, city, country string) (*Location, error) {
	if id == "" {
		return nil, errors.New("location ID is required")
	}

	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if city != "" {
		location.City = strings.TrimSpace(city)
	}

	if country != "" {
		location.Country = strings.TrimSpace(country)
	}

	if err := s.repo.Update(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

// DeleteLocation deletes a location
func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("location ID is required")
	}

	// Check if location is in use
	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if location.UsedByCount > 0 {
		return fmt.Errorf("cannot delete location with %d users", location.UsedByCount)
	}

	return s.repo.Delete(ctx, id)
}

// IncrementUsage increments the used_by_count
func (s *Service) IncrementUsage(ctx context.Context, id string) (int, error) {
	if id == "" {
		return 0, errors.New("location ID is required")
	}

	return s.repo.IncrementUsedByCount(ctx, id)
}

// DecrementUsage decrements the used_by_count
func (s *Service) DecrementUsage(ctx context.Context, id string) (int, error) {
	if id == "" {
		return 0, errors.New("location ID is required")
	}

	return s.repo.DecrementUsedByCount(ctx, id)
}
