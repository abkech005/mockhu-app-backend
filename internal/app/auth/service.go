package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles business logic for authentication operations.
// It coordinates between the HTTP layer and the data layer (repository).
type Service struct {
	repo UserRepository
}

// NewService creates a new authentication service instance.
// It requires a UserRepository implementation to interact with the database.
func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

// Signup creates a new user account with the provided information.
// It performs the following operations:
//   - Validates that the email doesn't already exist
//   - Hashes the password securely using bcrypt
//   - Generates a unique UUID for the user
//   - Creates the user record in the database
//
// Returns the created user (without password) or an error if the operation fails.
func (s *Service) Signup(ctx context.Context, email, password, firstName, lastName string) (*User, error) {
	// Check if user already exists
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user with their email and password.
// It performs the following validations:
//   - Checks if the user exists
//   - Verifies the account is active
//   - Validates the password against the stored hash
//   - Updates the last login timestamp
//
// Returns the authenticated user or an error if authentication fails.
func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login timestamp
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return user, nil
}

// GetUserByID retrieves a user by their unique identifier.
// Returns the user if found, or an error if the user doesn't exist.
func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// GetUserByEmail retrieves a user by their email address.
// Returns the user if found, or an error if the user doesn't exist.
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}

// VerifyEmail marks a user's email as verified.
// This is typically called after the user clicks a verification link.
// Returns an error if the user doesn't exist or the operation fails.
func (s *Service) VerifyEmail(ctx context.Context, userID string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

// UpdateProfile updates a user's profile information.
// Only non-sensitive fields like first name, last name, and avatar can be updated.
// Returns an error if the user doesn't exist or the operation fails.
func (s *Service) UpdateProfile(ctx context.Context, userID, firstName, lastName, avatarURL string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

// ChangePassword updates a user's password after verifying the old password.
// It performs the following operations:
//   - Verifies the old password is correct
//   - Hashes the new password
//   - Updates the user record
//
// Returns an error if the old password is incorrect or the operation fails.
func (s *Service) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}
