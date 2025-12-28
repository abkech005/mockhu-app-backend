package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetAuthURL returns the authorization URL for the specified provider
func (s *Service) GetAuthURL(provider, state string) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", fmt.Errorf("provider %s not configured", provider)
	}
	return p.GetAuthURL(state), nil
}

// OAuthResult contains the result of an OAuth authentication
type OAuthResult struct {
	User      *User
	IsNewUser bool
	Token     *OAuthToken
}

// OAuthSignupOrLogin handles OAuth-based authentication
func (s *Service) OAuthSignupOrLogin(ctx context.Context, provider, code string) (*OAuthResult, error) {
	// 1. Get Access Token
	p, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}

	token, err := p.ExchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// 2. Get User Info
	var userInfo *OAuthUserInfo

	// Special handling for Apple: if IDToken is present, try to use it
	if provider == "apple" && token.IDToken != "" {
		if appleProvider, ok := p.(*AppleOAuthProvider); ok {
			sub, email, err := appleProvider.ValidateIDToken(token.IDToken)
			if err == nil {
				userInfo = &OAuthUserInfo{
					ProviderID: sub,
					Email:      email,
				}
			}
		}
	}

	if userInfo == nil {
		userInfo, err = p.GetUserInfo(token.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}
	}

	// 3. Check if OAuth account exists
	oauthAccount, err := s.oauthRepo.FindByProvider(ctx, provider, userInfo.ProviderID)

	var user *User
	isNewUser := false

	if err == nil {
		// OAuth account exists -> Get User
		user, err = s.repo.FindByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user for oauth account: %w", err)
		}
	} else {
		// OAuth account doesn't exist -> Check if user with same email exists
		if userInfo.Email != "" {
			user, err = s.repo.FindByEmail(ctx, userInfo.Email)
		}

		if user != nil {
			// User exists -> Link Account
			err = s.linkOAuthAccount(ctx, user.ID, provider, userInfo, token)
			if err != nil {
				return nil, err
			}
		} else {
			// New User -> Create User & Link
			isNewUser = true

			// Generate random password
			randomPass := uuid.New().String()
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPass), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}

			user = &User{
				ID:            uuid.New().String(),
				Email:         userInfo.Email,
				EmailVerified: userInfo.EmailVerified,
				FirstName:     userInfo.FirstName,
				LastName:      userInfo.LastName,
				AvatarURL:     userInfo.AvatarURL,
				PasswordHash:  string(hashedPassword),
				IsActive:      true,
				// Set defaults
				WhoCanMessage:     "everyone",
				WhoCanSeePosts:    "everyone",
				ShowFollowersList: true,
				ShowFollowingList: true,
			}

			// If email is missing, generate a placeholder
			if user.Email == "" {
				user.Email = fmt.Sprintf("%s+%s@oauth.temp", provider, userInfo.ProviderID)
			}

			// Create user
			err = s.repo.Create(ctx, user)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			// Link OAuth
			err = s.linkOAuthAccount(ctx, user.ID, provider, userInfo, token)
			if err != nil {
				return nil, err
			}
		}
	}

	// Update tokens if existing user (simplified: we rely on initial link for now,
	// typically we'd update tokens here if we had Update method)

	return &OAuthResult{
		User:      user,
		IsNewUser: isNewUser,
		Token:     token,
	}, nil
}

// LinkOAuthProvider links an OAuth provider to an existing user
func (s *Service) LinkOAuthProvider(ctx context.Context, userID, provider, code string) (*OAuthUserInfo, error) {
	p, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}

	token, err := p.ExchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get User Info
	var userInfo *OAuthUserInfo
	if provider == "apple" && token.IDToken != "" {
		if appleProvider, ok := p.(*AppleOAuthProvider); ok {
			sub, email, err := appleProvider.ValidateIDToken(token.IDToken)
			if err == nil {
				userInfo = &OAuthUserInfo{
					ProviderID: sub,
					Email:      email,
				}
			}
		}
	}

	if userInfo == nil {
		userInfo, err = p.GetUserInfo(token.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}
	}

	// Check if already linked
	existing, err := s.oauthRepo.FindByProvider(ctx, provider, userInfo.ProviderID)
	if err == nil && existing.UserID != userID {
		return nil, fmt.Errorf("this %s account is already linked to another user", provider)
	}
	if err == nil {
		return userInfo, nil // Already linked to this user
	}

	err = s.linkOAuthAccount(ctx, userID, provider, userInfo, token)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// UnlinkOAuthProvider removes an OAuth provider link
func (s *Service) UnlinkOAuthProvider(ctx context.Context, userID, provider string) error {
	return s.oauthRepo.DeleteByUserID(ctx, userID, provider)
}

// linkOAuthAccount helper
func (s *Service) linkOAuthAccount(ctx context.Context, userID, provider string, userInfo *OAuthUserInfo, token *OAuthToken) error {
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	account := &OAuthAccount{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderID,
		ProviderEmail:  userInfo.Email,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: &expiresAt,
	}

	return s.oauthRepo.Create(ctx, account)
}
