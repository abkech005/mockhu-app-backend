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
		// Return actual error for debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get profile",
			"details": err.Error(),
		})
	}

	return c.JSON(profile)
}

// UpdateProfile handles PUT /v1/users/me/profile
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse request body
	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Update profile
	profile, err := h.service.UpdateProfile(c.Context(), currentUserID, &req)
	if err != nil {
		// Handle specific errors
		errMsg := err.Error()
		if errMsg == "username already taken" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username already taken",
			})
		}
		if errMsg == "no fields to update" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "no fields to update",
			})
		}
		// Validation errors
		if errMsg == "first name must be between 1 and 50 characters" ||
			errMsg == "last name must be between 1 and 50 characters" ||
			errMsg == "bio must not exceed 500 characters" ||
			errMsg == "username must be between 3 and 30 characters" ||
			errMsg == "username can only contain letters, numbers, and underscores" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update profile",
		})
	}

	return c.JSON(profile)
}

// UploadAvatar handles POST /v1/users/me/avatar
func (h *Handler) UploadAvatar(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse multipart form
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "avatar file is required",
		})
	}

	// Open and read file
	fileHandle, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read file",
		})
	}
	defer fileHandle.Close()

	// Read file bytes
	fileBytes := make([]byte, file.Size)
	_, err = fileHandle.Read(fileBytes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read file",
		})
	}

	// Upload avatar
	response, err := h.service.UploadAvatar(c.Context(), currentUserID, fileBytes, file.Filename)
	if err != nil {
		// Handle specific errors
		errMsg := err.Error()
		if errMsg == "file size exceeds 5MB" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "file size exceeds 5MB",
			})
		}
		if errMsg == "invalid file type, only JPEG, PNG, and WebP allowed" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid file type, only JPEG, PNG, and WebP allowed",
			})
		}
		if errMsg == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload avatar",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteAvatar handles DELETE /v1/users/me/avatar
func (h *Handler) DeleteAvatar(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Delete avatar
	err := h.service.DeleteAvatar(c.Context(), currentUserID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete avatar",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "avatar deleted successfully",
	})
}

// GetPrivacySettings handles GET /v1/users/me/privacy
func (h *Handler) GetPrivacySettings(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get privacy settings
	settings, err := h.service.GetPrivacySettings(c.Context(), currentUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get privacy settings",
		})
	}

	return c.JSON(settings)
}

// UpdatePrivacySettings handles PUT /v1/users/me/privacy
func (h *Handler) UpdatePrivacySettings(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse request body
	var req UpdatePrivacyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Update privacy settings
	settings, err := h.service.UpdatePrivacySettings(c.Context(), currentUserID, &req)
	if err != nil {
		// Handle validation errors
		errMsg := err.Error()
		if errMsg == "who_can_message must be 'everyone', 'followers', or 'none'" ||
			errMsg == "who_can_see_posts must be 'everyone', 'followers', or 'none'" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update privacy settings",
		})
	}

	return c.JSON(settings)
}

// GetMutualConnections handles GET /v1/users/:userId/mutual-connections
func (h *Handler) GetMutualConnections(c *fiber.Ctx) error {
	// Get current user ID from JWT (required)
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get target user ID from URL params
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	// Get mutual connections
	response, err := h.service.GetMutualConnections(c.Context(), currentUserID, targetUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get mutual connections",
		})
	}

	return c.JSON(response)
}

// GetMutualConnectionsCount handles GET /v1/users/:userId/mutual-connections/count
func (h *Handler) GetMutualConnectionsCount(c *fiber.Ctx) error {
	// Get current user ID from JWT (required)
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get target user ID from URL params
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Get mutual connections count
	count, err := h.service.GetMutualConnectionsCount(c.Context(), currentUserID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get mutual connections count",
		})
	}

	return c.JSON(MutualConnectionsCountResponse{
		UserID:                 targetUserID,
		MutualConnectionsCount: count,
	})
}
