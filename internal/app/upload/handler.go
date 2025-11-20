package upload

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	// Add storage service later
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Avatar(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file provided"})
	}

	_ = file // Use file later

	return c.JSON(AvatarResponse{
		AvatarURL: "https://cdn.example.com/avatars/uploaded.jpg",
	})
}
