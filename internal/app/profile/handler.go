package profile

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for profile operations
type Handler struct {
	service ProfileService
}

// NewHandler creates a new profile handler
func NewHandler(service ProfileService) *Handler {
	return &Handler{
		service: service,
	}
}

// Handlers will be implemented in Phase 3+

// GetUserProfile handles GET /v1/users/:userId/profile
func (h *Handler) GetUserProfile(c *fiber.Ctx) error {
	// Get user ID from URL params
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Get current user ID from context (optional - for authenticated users)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get profile
	profile, err := h.service.GetUserProfile(c.Context(), userID, currentUserID)
	if err != nil {
		// Check if user not found
		if err.Error() == "user not found" || err.Error() == "user not found: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get profile",
		})
	}

	return c.JSON(profile)
}

// GetOwnProfile handles GET /v1/users/me/profile
func (h *Handler) GetOwnProfile(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get profile
	profile, err := h.service.GetOwnProfile(c.Context(), currentUserID)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "user not found: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get profile",
		})
	}

	return c.JSON(profile)
}

// UpdateProfile handles PUT /v1/users/me/profile
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	// TODO: Implement in Phase 4
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// UploadAvatar handles POST /v1/users/me/avatar
func (h *Handler) UploadAvatar(c *fiber.Ctx) error {
	// TODO: Implement in Phase 5
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// DeleteAvatar handles DELETE /v1/users/me/avatar
func (h *Handler) DeleteAvatar(c *fiber.Ctx) error {
	// TODO: Implement in Phase 5
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// GetPrivacySettings handles GET /v1/users/me/privacy
func (h *Handler) GetPrivacySettings(c *fiber.Ctx) error {
	// TODO: Implement in Phase 6
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// UpdatePrivacySettings handles PUT /v1/users/me/privacy
func (h *Handler) UpdatePrivacySettings(c *fiber.Ctx) error {
	// TODO: Implement in Phase 6
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// GetMutualConnections handles GET /v1/users/:userId/mutual-connections
func (h *Handler) GetMutualConnections(c *fiber.Ctx) error {
	// TODO: Implement in Phase 7
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}

// GetMutualConnectionsCount handles GET /v1/users/:userId/mutual-connections/count
func (h *Handler) GetMutualConnectionsCount(c *fiber.Ctx) error {
	// TODO: Implement in Phase 7
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "not implemented yet",
	})
}
