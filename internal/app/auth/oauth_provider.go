package auth

// OAuthProvider defines the common interface for all OAuth providers
type OAuthProvider interface {
	// GetAuthURL returns the URL to redirect the user to for authentication
	GetAuthURL(state string) string

	// ExchangeCode exchanges the authorization code for tokens
	ExchangeCode(code string) (*OAuthToken, error)

	// GetUserInfo retrieves the user's information using the access token
	GetUserInfo(token string) (*OAuthUserInfo, error)
}

// OAuthToken represents the tokens returned by the provider
type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
	IDToken      string // Specific to OIDC providers like Google/Apple
}

// OAuthUserInfo represents the normalized user information
type OAuthUserInfo struct {
	ProviderID    string
	Email         string
	EmailVerified bool
	FirstName     string
	LastName      string
	AvatarURL     string
}
