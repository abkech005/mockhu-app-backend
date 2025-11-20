package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	// Add dependencies here later
	// service *Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(SignupResponse{
		UserID:              "user-uuid-12345",
		VerificationNeeded:  true,
		VerificationChannel: req.Method,
	})
}

func (h *Handler) Verify(c *fiber.Ctx) error {
	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

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

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(LoginResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy",
		RefreshToken: "refresh-token-dummy",
		ExpiresIn:    900,
	})
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(RefreshResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new-dummy",
		RefreshToken: "new-refresh-token-dummy",
		ExpiresIn:    900,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(LogoutResponse{
		Message: "logged_out",
	})
}

func (h *Handler) Resend(c *fiber.Ctx) error {
	var req ResendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ResendResponse{
		Message: "code_sent",
	})
}
