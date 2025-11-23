package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles business logic for authentication operations.
// It coordinates between the HTTP layer and the data layer (repository).
type Service struct {
	repo             UserRepository
	verificationRepo VerificationRepository
}

// NewService creates a new authentication service instance.
// It requires a UserRepository and VerificationRepository to interact with the database.
func NewService(repo UserRepository, verificationRepo VerificationRepository) *Service {
	return &Service{
		repo:             repo,
		verificationRepo: verificationRepo,
	}
}

// SignupResult contains the user and optional verification code
type SignupResult struct {
	User              *User
	VerificationCode  *VerificationCode
	NeedsVerification bool
}

// Signup creates a new user account with the provided information.
// It performs the following operations:
//   - Validates that the email/phone doesn't already exist
//   - Hashes the password securely using bcrypt
//   - Generates a unique UUID for the user
//   - Creates the user record in the database
//   - Automatically sends verification code based on signup method
//
// Returns the created user, verification code (if applicable), or an error if the operation fails.
func (s *Service) Signup(ctx context.Context, method, email, phone, password string) (*SignupResult, error) {
	// Validate based on method
	if method == "email" && email == "" {
		return nil, errors.New("email is required for email signup")
	}
	if method == "mobile" && phone == "" {
		return nil, errors.New("phone is required for mobile signup")
	}

	// Check if user already exists
	if email != "" {
		existing, _ := s.repo.FindByEmail(ctx, email)
		if existing != nil {
			return nil, errors.New("email already registered")
		}
	}
	if phone != "" {
		existing, _ := s.repo.FindByPhone(ctx, phone)
		if existing != nil {
			return nil, errors.New("phone already registered")
		}
	}

	// Hash password (if provided)
	var hashedPassword string
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		hashedPassword = string(hashed)
	}

	// Create user
	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	result := &SignupResult{
		User:              user,
		NeedsVerification: false,
	}

	// Auto-send verification based on signup method
	switch method {
	case "email":
		verificationCode, err := s.GenerateEmailVerificationCode(ctx, user.ID)
		if err != nil {
			log.Printf("⚠️ Failed to send email verification: %v", err)
			// Don't fail signup if verification fails
		} else {
			result.VerificationCode = verificationCode
			result.NeedsVerification = true
		}

	case "mobile":
		verificationCode, err := s.GeneratePhoneVerificationCode(ctx, user.ID, phone)
		if err != nil {
			log.Printf("⚠️ Failed to send phone verification: %v", err)
			// Don't fail signup if verification fails
		} else {
			result.VerificationCode = verificationCode
			result.NeedsVerification = true
		}

	case "google", "facebook":
		// Social signups are pre-verified
		user.EmailVerified = true
		user.UpdatedAt = time.Now()
		_ = s.repo.Update(ctx, user)
	}

	return result, nil
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

// GenerateEmailVerificationCode creates a new 6-digit verification code for email verification.
// It deactivates any previous active codes and logs the code (since email infrastructure isn't implemented yet).
func (s *Service) GenerateEmailVerificationCode(ctx context.Context, userID string) (*VerificationCode, error) {
	// Get user to retrieve email
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Deactivate previous codes for this user (optimal: just sets is_active=false)
	_ = s.verificationRepo.DeactivatePreviousCodes(ctx, userID, VerificationTypeEmail)

	// Generate 6-digit code
	code := generateRandomCode()

	verification := &VerificationCode{
		ID:        uuid.New().String(),
		UserID:    userID,
		Code:      code,
		Type:      VerificationTypeEmail,
		Contact:   user.Email,
		IsActive:  true, // New codes are active by default
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to create verification code: %w", err)
	}

	// TODO: Send email when notification infrastructure is ready
	log.Printf("📧 [MOCK EMAIL] To: %s | Code: %s | Expires in 10 minutes", user.Email, code)

	return verification, nil
}

// VerifyEmailCode validates an email verification code and marks the user's email as verified.
// It checks that the code is active, exists, hasn't been used, and hasn't expired.
func (s *Service) VerifyEmailCode(ctx context.Context, userID, code string) error {
	// Find the verification code (only returns active codes)
	verification, err := s.verificationRepo.FindByCodeAndType(ctx, code, VerificationTypeEmail)
	if err != nil {
		return errors.New("invalid or expired code")
	}

	// Verify it belongs to the user
	if verification.UserID != userID {
		return errors.New("invalid code")
	}

	// Check if active (should be true from query, but double-check)
	if !verification.IsActive {
		return errors.New("code is no longer active")
	}

	// Check if already used
	if verification.UsedAt != nil {
		return errors.New("code already used")
	}

	// Check if expired
	if time.Now().After(verification.ExpiresAt) {
		return errors.New("code has expired")
	}

	// Mark code as used (also sets is_active=false)
	if err := s.verificationRepo.MarkAsUsed(ctx, verification.ID); err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	// Update user's email_verified status
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("✅ Email verified for user: %s", userID)
	return nil
}

// GeneratePhoneVerificationCode creates a new 6-digit verification code for phone verification.
// It deactivates any previous active codes and logs the code (since SMS infrastructure isn't implemented yet).
func (s *Service) GeneratePhoneVerificationCode(ctx context.Context, userID, phoneNumber string) (*VerificationCode, error) {
	// Verify user exists
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update user's phone number if provided
	if phoneNumber != "" && user.Phone != phoneNumber {
		user.Phone = phoneNumber
		user.UpdatedAt = time.Now()
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update phone number: %w", err)
		}
	}

	// Deactivate previous codes for this user (optimal: just sets is_active=false)
	_ = s.verificationRepo.DeactivatePreviousCodes(ctx, userID, VerificationTypePhone)

	// Generate 6-digit code
	code := generateRandomCode()

	verification := &VerificationCode{
		ID:        uuid.New().String(),
		UserID:    userID,
		Code:      code,
		Type:      VerificationTypePhone,
		Contact:   user.Phone,
		IsActive:  true, // New codes are active by default
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to create verification code: %w", err)
	}

	// TODO: Send SMS when notification infrastructure is ready
	log.Printf("📱 [MOCK SMS] To: %s | Code: %s | Expires in 10 minutes", user.Phone, code)

	return verification, nil
}

// VerifyPhoneCode validates a phone verification code and marks the user's phone as verified.
// It checks that the code is active, exists, hasn't been used, and hasn't expired.
func (s *Service) VerifyPhoneCode(ctx context.Context, userID, code string) error {
	// Find the verification code (only returns active codes)
	verification, err := s.verificationRepo.FindByCodeAndType(ctx, code, VerificationTypePhone)
	if err != nil {
		return errors.New("invalid or expired code")
	}

	// Verify it belongs to the user
	if verification.UserID != userID {
		return errors.New("invalid code")
	}

	// Check if active (should be true from query, but double-check)
	if !verification.IsActive {
		return errors.New("code is no longer active")
	}

	// Check if already used
	if verification.UsedAt != nil {
		return errors.New("code already used")
	}

	// Check if expired
	if time.Now().After(verification.ExpiresAt) {
		return errors.New("code has expired")
	}

	// Mark code as used (also sets is_active=false)
	if err := s.verificationRepo.MarkAsUsed(ctx, verification.ID); err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	// Update user's phone_verified status
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.PhoneVerified = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("✅ Phone verified for user: %s", userID)
	return nil
}

// generateRandomCode generates a secure random 6-digit verification code.
func generateRandomCode() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to timestamp-based code if crypto/rand fails
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}
