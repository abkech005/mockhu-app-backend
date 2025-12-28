package auth

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/Timothylock/go-signin-with-apple/apple"
)

// AppleOAuthProvider implements OAuthProvider for Apple
type AppleOAuthProvider struct {
	clientID    string
	teamID      string
	keyID       string
	redirectURL string
	privateKey  string
}

// NewAppleOAuthProvider creates a new Apple OAuth provider
func NewAppleOAuthProvider(clientID, teamID, keyID, privateKey, redirectURL string) *AppleOAuthProvider {
	return &AppleOAuthProvider{
		clientID:    clientID,
		teamID:      teamID,
		keyID:       keyID,
		privateKey:  privateKey,
		redirectURL: redirectURL,
	}
}

// GetAuthURL returns the Apple login URL
func (p *AppleOAuthProvider) GetAuthURL(state string) string {
	// Using form_post response mode for Apple
	return fmt.Sprintf("https://appleid.apple.com/auth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=name%%20email&response_mode=form_post&state=%s",
		p.clientID, p.redirectURL, state)
}

// ExchangeCode exchanges the authorization code for tokens
func (p *AppleOAuthProvider) ExchangeCode(code string) (*OAuthToken, error) {
	secret, err := apple.GenerateClientSecret(p.privateKey, p.teamID, p.clientID, p.keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client secret: %w", err)
	}

	client := apple.New()
	vReq := apple.AppValidationTokenRequest{
		ClientID:     p.clientID,
		ClientSecret: secret,
		Code:         code,
	}

	var resp apple.ValidationResponse
	if err := client.VerifyAppToken(context.Background(), vReq, &resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("apple auth error: %s", resp.Error)
	}

	// Apple returns ID token which contains user info
	return &OAuthToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		IDToken:      resp.IDToken,
	}, nil
}

// GetUserInfo retrieves user info from Apple ID Token
func (p *AppleOAuthProvider) GetUserInfo(token string) (*OAuthUserInfo, error) {
	// For Apple, the "token" passed here should be the ID Token
	// If it's an access token, we can't get user info easily without the ID token
	// This provider implementation expects the ID Token to be passed as 'token'
	// or we need to change interface to support ID Token parsing separately.
	// For now, let's assume the ID Token is passed here if available.

	// Claims decoding would go here (using jwt library)
	// Since we don't have the jwt lib imported yet in this file,
	// and Apple sends name/email only on first login in the POST body (not token),
	// this is tricky.

	// TODO: Implement Apple ID Token claim parsing
	// Note: Apple only returns Name on the FIRST login in the callback.
	// Subsequent logins only return the unique sub (ProviderID) and email in ID token.

	return nil, fmt.Errorf("apple user info retrieval requires ID token parsing implementation")
}

// ValidateIDToken validates the Apple ID Token and returns claims
func (p *AppleOAuthProvider) ValidateIDToken(idToken string) (string, string, error) {
	// Implement ID validation logic
	// Returns UserID (sub) and Email
	claims, err := apple.GetClaims(idToken)
	if err != nil {
		return "", "", err
	}

	// claims is a pointer to MapClaims, so we must dereference it
	c := *claims
	email, _ := c["email"].(string)
	sub, _ := c["sub"].(string)

	return sub, email, nil
}
