package share

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for share operations
type Handler struct {
	service ShareService
}

// NewHandler creates a new share handler
func NewHandler(service ShareService) *Handler {
	return &Handler{service: service}
}

// CreateShare handles POST /v1/posts/:postId/shares
func (h *Handler) CreateShare(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get post ID from URL
	postID := c.Params("postId")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post ID is required",
		})
	}

	// Parse request body
	var req CreateShareRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if req.SharedToType == "" {
		req.SharedToType = "timeline" // Default to timeline
	}

	if req.SharedToType != "timeline" && req.SharedToType != "dm" && req.SharedToType != "external" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid share type. Must be 'timeline', 'dm', or 'external'",
		})
	}

	// Create share
	share, err := h.service.CreateShare(c.Context(), postID, currentUserID, &req)
	if err != nil {
		if err == ErrInvalidShareType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid share type",
			})
		}
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post not found",
			})
		}
		if err == ErrAlreadyShared {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "post already shared by user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create share",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(share)
}

// GetShare handles GET /v1/shares/:shareId
func (h *Handler) GetShare(c *fiber.Ctx) error {
	// Get share ID from URL
	shareID := c.Params("shareId")
	if shareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "share ID is required",
		})
	}

	// Get current user ID (optional)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get share
	share, err := h.service.GetShare(c.Context(), shareID, currentUserID)
	if err != nil {
		if err == ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get share",
		})
	}

	return c.JSON(share)
}

// GetPostShares handles GET /v1/posts/:postId/shares
func (h *Handler) GetPostShares(c *fiber.Ctx) error {
	// Get post ID from URL
	postID := c.Params("postId")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post ID is required",
		})
	}

	// Get current user ID (optional)
	currentUserID, _ := c.Locals("user_id").(string)

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get shares
	response, err := h.service.GetPostShares(c.Context(), postID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get shares",
		})
	}

	return c.JSON(response)
}

// GetUserShares handles GET /v1/users/:userId/shares
func (h *Handler) GetUserShares(c *fiber.Ctx) error {
	// Get user ID from URL
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Get current user ID (optional)
	currentUserID, _ := c.Locals("user_id").(string)

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get shares
	response, err := h.service.GetUserShares(c.Context(), userID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get shares",
		})
	}

	return c.JSON(response)
}

// DeleteShare handles DELETE /v1/shares/:shareId
func (h *Handler) DeleteShare(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get share ID from URL
	shareID := c.Params("shareId")
	if shareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "share ID is required",
		})
	}

	// Delete share
	err := h.service.DeleteShare(c.Context(), shareID, currentUserID)
	if err != nil {
		if err == ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "share not found",
			})
		}
		if err == ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "unauthorized to delete this share",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete share",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "share deleted successfully",
	})
}

// GetShareCount handles GET /v1/posts/:postId/shares/count
func (h *Handler) GetShareCount(c *fiber.Ctx) error {
	// Get post ID from URL
	postID := c.Params("postId")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post ID is required",
		})
	}

	// Get share count
	count, err := h.service.GetShareCount(c.Context(), postID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get share count",
		})
	}

	return c.JSON(fiber.Map{
		"post_id": postID,
		"count":   count,
	})
}


