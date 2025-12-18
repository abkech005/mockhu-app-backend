package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
)

// OAuthInitiate redirects user to OAuth provider
func (h *Handler) OAuthInitiate(c *fiber.Ctx) error {
	provider := c.Params("provider")

	// Generate random state
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in cookie (valid for 10 mins)
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   true, // Should check if using HTTPS
		SameSite: "Lax",
	})

	url, err := h.service.GetAuthURL(provider, state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Redirect(url)
}

// OAuthCallback handles OAuth provider callback
func (h *Handler) OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")

	// Verify state
	cookieState := c.Cookies("oauth_state")
	if state != cookieState {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid oauth state",
		})
	}

	// Clear state cookie
	c.ClearCookie("oauth_state")

	result, err := h.service.OAuthSignupOrLogin(c.Context(), provider, code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// In a real app, you would likely generate your OWN JWT tokens here
	// using result.User.ID.
	// Since login/signup logic in Handler usually does this, we should reuse it.
	// But `OAuthSignupOrLogin` returned user and we need tokens.
	// We haven't implemented JWT generation helper in Service yet (it was TODO).
	// But `Login` in handler generates dummy tokens. I should duplicate that logic or extract it.

	// Duplicate dummy token logic for now as per handler.go
	accessToken := "dummy_access_token_" + result.User.ID
	refreshToken := "dummy_refresh_token_" + result.User.ID

	return c.JSON(OAuthCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		IsNewUser:    result.IsNewUser,
		User: &UserInfo{
			ID:       result.User.ID,
			Username: result.User.Username,
			Email:    result.User.Email,
		},
	})
}

// OAuthLink links OAuth provider to existing account
func (h *Handler) OAuthLink(c *fiber.Ctx) error {
	provider := c.Params("provider")
	var req OAuthLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Helper to get userID from context (set by AuthMiddleware)
	// assuming "user_id" is set in locals
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userInfo, err := h.service.LinkOAuthProvider(c.Context(), userID, provider, req.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(OAuthLinkResponse{
		Message:       "Account linked successfully",
		Provider:      provider,
		ProviderEmail: userInfo.Email,
	})
}

// OAuthUnlink removes OAuth provider link
func (h *Handler) OAuthUnlink(c *fiber.Ctx) error {
	provider := c.Params("provider")

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err := h.service.UnlinkOAuthProvider(c.Context(), userID, provider)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account unlinked successfully",
	})
}
