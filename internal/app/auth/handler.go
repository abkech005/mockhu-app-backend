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

// SendEmailVerification generates and sends an email verification code.
// The code is logged to console since email infrastructure isn't implemented yet.
func (h *Handler) SendEmailVerification(c *fiber.Ctx) error {
	var req SendEmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	verification, err := h.service.GenerateEmailVerificationCode(c.Context(), req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(SendEmailVerificationResponse{
		Message:   "Verification code sent to your email",
		ExpiresIn: 600,               // 10 minutes
		Code:      verification.Code, // TODO: Remove in production (for testing only)
	})
}

// VerifyEmail validates an email verification code and marks the user's email as verified.
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.service.VerifyEmailCode(c.Context(), req.UserID, req.Code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(VerifyEmailResponse{
		Message:       "Email verified successfully",
		EmailVerified: true,
	})
}

// SendPhoneVerification generates and sends a phone verification code.
// The code is logged to console since SMS infrastructure isn't implemented yet.
func (h *Handler) SendPhoneVerification(c *fiber.Ctx) error {
	var req SendPhoneVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	verification, err := h.service.GeneratePhoneVerificationCode(c.Context(), req.UserID, req.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(SendPhoneVerificationResponse{
		Message:   "Verification code sent to your phone",
		ExpiresIn: 600,               // 10 minutes
		Code:      verification.Code, // TODO: Remove in production (for testing only)
	})
}

// VerifyPhone validates a phone verification code and marks the user's phone as verified.
func (h *Handler) VerifyPhone(c *fiber.Ctx) error {
	var req VerifyPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.service.VerifyPhoneCode(c.Context(), req.UserID, req.Code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(VerifyPhoneResponse{
		Message:       "Phone verified successfully",
		PhoneVerified: true,
	})
}
