package onboarding

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	// Add dependencies here later
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Basic(c *fiber.Ctx) error {
	var req BasicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(BasicResponse{
		TempID:  "temp-uuid-67890",
		Message: "basic_saved",
	})
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	var req ProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ProfileResponse{
		TempID:            req.TempID,
		UsernameAvailable: true,
		AvatarURL:         "https://cdn.example.com/avatars/dummy.jpg",
	})
}

func (h *Handler) Interests(c *fiber.Ctx) error {
	var req InterestsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(InterestsResponse{
		Message: "interests_saved",
	})
}
