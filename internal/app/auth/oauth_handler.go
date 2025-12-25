package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"mockhu-app-backend/internal/pkg/jwt"

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

	// Generate real JWT tokens
	accessToken, err := jwt.GenerateAccessToken(result.User.ID, result.User.Email, result.User.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate access token",
		})
	}

	refreshToken, err := jwt.GenerateRefreshToken(result.User.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate refresh token",
		})
	}

	return c.JSON(OAuthCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(jwt.AccessTokenDuration.Seconds()),
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
