package follow

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for follow operations
type Handler struct {
	service FollowService
}

// NewHandler creates a new follow handler
func NewHandler(service FollowService) *Handler {
	return &Handler{service: service}
}

// Follow handles POST /v1/users/:userId/follow
func (h *Handler) Follow(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Follow user
	result, err := h.service.Follow(c.Context(), currentUserID, targetUserID)
	if err != nil {
		if err == ErrCannotFollowSelf {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot follow yourself",
			})
		}
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to follow user",
		})
	}

	return c.JSON(result)
}

// Unfollow handles DELETE /v1/users/:userId/follow
func (h *Handler) Unfollow(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Unfollow user
	result, err := h.service.Unfollow(c.Context(), currentUserID, targetUserID)
	if err != nil {
		if err == ErrCannotFollowSelf {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot unfollow yourself",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unfollow user",
		})
	}

	return c.JSON(result)
}

// GetFollowers handles GET /v1/users/:userId/followers
func (h *Handler) GetFollowers(c *fiber.Ctx) error {
	// Get current user ID from JWT (optional for checking if we follow them)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Parse pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get followers
	result, err := h.service.GetFollowers(c.Context(), targetUserID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get followers",
		})
	}

	return c.JSON(result)
}

// GetFollowing handles GET /v1/users/:userId/following
func (h *Handler) GetFollowing(c *fiber.Ctx) error {
	// Get current user ID from JWT (optional for checking if we follow them)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Parse pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get following
	result, err := h.service.GetFollowing(c.Context(), targetUserID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get following",
		})
	}

	return c.JSON(result)
}

// IsFollowing handles GET /v1/users/:userId/is-following
func (h *Handler) IsFollowing(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Check if following
	result, err := h.service.IsFollowing(c.Context(), currentUserID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check follow status",
		})
	}

	return c.JSON(result)
}

// GetFollowStats handles GET /v1/users/:userId/follow-stats
func (h *Handler) GetFollowStats(c *fiber.Ctx) error {
	// Get target user ID from URL
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Get stats
	result, err := h.service.GetFollowStats(c.Context(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get follow stats",
		})
	}

	return c.JSON(result)
}
