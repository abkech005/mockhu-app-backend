package auth

import (
	"context"
	"time"

	fb "github.com/huandu/facebook/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// FacebookOAuthProvider implements OAuthProvider for Facebook
type FacebookOAuthProvider struct {
	config *oauth2.Config
}

// NewFacebookOAuthProvider creates a new Facebook OAuth provider
func NewFacebookOAuthProvider(clientID, clientSecret, redirectURL string) *FacebookOAuthProvider {
	return &FacebookOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "public_profile"},
			Endpoint:     facebook.Endpoint,
		},
	}
}

// GetAuthURL returns the Facebook login URL
func (p *FacebookOAuthProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

// ExchangeCode exchanges the authorization code for tokens
func (p *FacebookOAuthProvider) ExchangeCode(code string) (*OAuthToken, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return &OAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Sub(time.Now()).Seconds()),
		TokenType:    token.TokenType,
	}, nil
}

// GetUserInfo retrieves user info from Facebook
func (p *FacebookOAuthProvider) GetUserInfo(token string) (*OAuthUserInfo, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "id,name,email,first_name,last_name,picture.width(400).height(400)",
		"access_token": token,
	})
	if err != nil {
		return nil, err
	}

	var email string
	if val, ok := res["email"].(string); ok {
		email = val
	}

	var firstName string
	if val, ok := res["first_name"].(string); ok {
		firstName = val
	}

	var lastName string
	if val, ok := res["last_name"].(string); ok {
		lastName = val
	}

	var avatarURL string
	if picture, ok := res["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			if url, ok := data["url"].(string); ok {
				avatarURL = url
			}
		}
	}

	return &OAuthUserInfo{
		ProviderID:    res["id"].(string),
		Email:         email,
		EmailVerified: email != "", // Facebook verifies emails
		FirstName:     firstName,
		LastName:      lastName,
		AvatarURL:     avatarURL,
	}, nil
}
