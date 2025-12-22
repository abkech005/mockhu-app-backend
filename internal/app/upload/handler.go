package upload

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RequestAvatarUpload generates a presigned PUT URL for direct upload.
// POST /v1/upload/avatar/request
func (h *Handler) RequestAvatarUpload(c *fiber.Ctx) error {
	var req AvatarUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	// Default content type
	if req.ContentType == "" {
		req.ContentType = "image/jpeg"
	}

	response, err := h.service.RequestAvatarUpload(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

// ConfirmAvatarUpload saves the avatar URL to user profile.
// POST /v1/upload/avatar/confirm
func (h *Handler) ConfirmAvatarUpload(c *fiber.Ctx) error {
	var req AvatarConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.UserID == "" || req.FileKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id and file_key are required",
		})
	}

	response, err := h.service.ConfirmAvatarUpload(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
