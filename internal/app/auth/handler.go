package auth

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for authentication endpoints.
// It uses the Service layer to perform business operations.
type Handler struct {
	service *Service
}

// NewHandler creates a new authentication handler with the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Signup handles new user registration.
// It validates the request, creates a new user, and returns the user ID.
// TODO: Add email/phone verification logic.
func (h *Handler) Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Create user via service
	// Note: FirstName and LastName are empty for now as they're not in the signup request
	user, err := h.service.Signup(c.Context(), req.Email, req.Password, "", "")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(SignupResponse{
		UserID:              user.ID,
		VerificationNeeded:  true,
		VerificationChannel: req.Method,
	})
}

// Verify handles email/phone verification.
// TODO: Implement verification code validation and token generation.
func (h *Handler) Verify(c *fiber.Ctx) error {
	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// TODO: Verify the code and mark email/phone as verified
	return c.JSON(VerifyResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy",
		RefreshToken: "refresh-token-dummy",
		ExpiresIn:    900,
		User: &UserInfo{
			ID:       req.UserID,
			Username: "john123",
			Email:    "user@example.com",
		},
	})
}

// Login handles user authentication.
// It validates credentials and returns access tokens.
// TODO: Generate real JWT tokens instead of dummy tokens.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Authenticate user via service
	user, err := h.service.Login(c.Context(), req.Identifier, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO: Generate real JWT tokens
	return c.JSON(LoginResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy-" + user.ID,
		RefreshToken: "refresh-token-dummy-" + user.ID,
		ExpiresIn:    900,
	})
}

// Refresh handles token refresh requests.
// TODO: Implement JWT token refresh logic.
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// TODO: Validate refresh token and issue new tokens
	return c.JSON(RefreshResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new-dummy",
		RefreshToken: "new-refresh-token-dummy",
		ExpiresIn:    900,
	})
}

// Logout handles user logout requests.
// TODO: Implement token invalidation/blacklisting.
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// TODO: Invalidate refresh token
	return c.JSON(LogoutResponse{
		Message: "logged_out",
	})
}

// Resend handles resending verification codes.
// TODO: Implement code generation and sending logic.
func (h *Handler) Resend(c *fiber.Ctx) error {
	var req ResendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// TODO: Generate and send new verification code
	return c.JSON(ResendResponse{
		Message: "code_sent",
	})
}
