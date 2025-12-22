package education

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service handles education business logic
type Service struct {
	repo EducationRepository
}

// NewService creates a new education service
func NewService(repo EducationRepository) *Service {
	return &Service{repo: repo}
}

// GetUserEducation retrieves all education entries for a user
func (s *Service) GetUserEducation(ctx context.Context, userID string) ([]Education, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	educations, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education entries: %w", err)
	}

	if educations == nil {
		educations = []Education{}
	}

	return educations, nil
}

// GetEducationByID retrieves an education entry by its ID
func (s *Service) GetEducationByID(ctx context.Context, id string) (*Education, error) {
	if id == "" {
		return nil, errors.New("education ID is required")
	}

	return s.repo.FindByID(ctx, id)
}

// GetCurrentEducation retrieves the current education for a user
func (s *Service) GetCurrentEducation(ctx context.Context, userID string) (*Education, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.repo.FindCurrentEducation(ctx, userID)
}

// CreateEducation creates a new education entry for a user
func (s *Service) CreateEducation(ctx context.Context, userID string, req CreateEducationRequest) (*Education, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Validate school
	school := strings.TrimSpace(req.School)
	if school == "" {
		return nil, errors.New("school name is required")
	}

	if len(school) > 255 {
		return nil, errors.New("school name must be 255 characters or less")
	}

	// Validate year logic
	if req.StartYear != nil && req.EndYear != nil && *req.EndYear < *req.StartYear {
		return nil, errors.New("end_year cannot be before start_year")
	}

	// If current is true, end_year should be nil
	if req.Current && req.EndYear != nil {
		return nil, errors.New("current education should not have an end_year")
	}

	education := &Education{
		UserID:       userID,
		School:       school,
		Degree:       strings.TrimSpace(req.Degree),
		FieldOfStudy: strings.TrimSpace(req.FieldOfStudy),
		Location:     strings.TrimSpace(req.Location),
		StartYear:    req.StartYear,
		EndYear:      req.EndYear,
		Current:      req.Current,
		LogoURL:      strings.TrimSpace(req.LogoURL),
		Grade:        strings.TrimSpace(req.Grade),
		Activities:   strings.TrimSpace(req.Activities),
		Description:  strings.TrimSpace(req.Description),
	}

	if err := s.repo.Create(ctx, education); err != nil {
		return nil, fmt.Errorf("failed to create education entry: %w", err)
	}

	return education, nil
}

// UpdateEducation updates an existing education entry
func (s *Service) UpdateEducation(ctx context.Context, id string, req UpdateEducationRequest) (*Education, error) {
	if id == "" {
		return nil, errors.New("education ID is required")
	}

	// Get existing education entry
	education, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.School != "" {
		school := strings.TrimSpace(req.School)
		if len(school) > 255 {
			return nil, errors.New("school name must be 255 characters or less")
		}
		education.School = school
	}

	if req.Degree != "" {
		education.Degree = strings.TrimSpace(req.Degree)
	}

	if req.FieldOfStudy != "" {
		education.FieldOfStudy = strings.TrimSpace(req.FieldOfStudy)
	}

	if req.Location != "" {
		education.Location = strings.TrimSpace(req.Location)
	}

	if req.StartYear != nil {
		education.StartYear = req.StartYear
	}

	if req.EndYear != nil {
		education.EndYear = req.EndYear
	}

	if req.Current != nil {
		education.Current = *req.Current
	}

	if req.LogoURL != "" {
		education.LogoURL = strings.TrimSpace(req.LogoURL)
	}

	if req.Grade != "" {
		education.Grade = strings.TrimSpace(req.Grade)
	}

	if req.Activities != "" {
		education.Activities = strings.TrimSpace(req.Activities)
	}

	if req.Description != "" {
		education.Description = strings.TrimSpace(req.Description)
	}

	// Validate year logic
	if education.StartYear != nil && education.EndYear != nil && *education.EndYear < *education.StartYear {
		return nil, errors.New("end_year cannot be before start_year")
	}

	if err := s.repo.Update(ctx, education); err != nil {
		return nil, fmt.Errorf("failed to update education entry: %w", err)
	}

	return education, nil
}

// DeleteEducation deletes an education entry
func (s *Service) DeleteEducation(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("education ID is required")
	}

	return s.repo.Delete(ctx, id)
}

// CountUserEducation returns the count of education entries for a user
func (s *Service) CountUserEducation(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, errors.New("user ID is required")
	}

	return s.repo.CountByUserID(ctx, userID)
}
