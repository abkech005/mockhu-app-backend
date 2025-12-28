package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthProvider implements OAuthProvider for Google
type GoogleOAuthProvider struct {
	config *oauth2.Config
}

// NewGoogleOAuthProvider creates a new Google OAuth provider
func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) *GoogleOAuthProvider {
	return &GoogleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL returns the Google login URL
func (p *GoogleOAuthProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges the authorization code for tokens
func (p *GoogleOAuthProvider) ExchangeCode(code string) (*OAuthToken, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	idToken, _ := token.Extra("id_token").(string)

	return &OAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Sub(time.Now()).Seconds()),
		TokenType:    token.TokenType,
		IDToken:      idToken,
	}, nil
}

// GetUserInfo retrieves user info from Google
func (p *GoogleOAuthProvider) GetUserInfo(token string) (*OAuthUserInfo, error) {
	// Create an HTTP client with the token
	client := p.config.Client(context.Background(), &oauth2.Token{AccessToken: token})

	// Request user info
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	// Parse response
	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		ProviderID:    googleUser.ID,
		Email:         googleUser.Email,
		EmailVerified: googleUser.VerifiedEmail,
		FirstName:     googleUser.GivenName,
		LastName:      googleUser.FamilyName,
		AvatarURL:     googleUser.Picture,
	}, nil
}
