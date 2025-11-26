package post

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for post operations
type Handler struct {
	service PostService
}

// NewHandler creates a new post handler
func NewHandler(service PostService) *Handler {
	return &Handler{service: service}
}

// CreatePost handles POST /v1/posts
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse request body
	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content is required",
		})
	}

	if len(req.Content) > 5000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content too long (max 5000 characters)",
		})
	}

	if len(req.Images) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "too many images (max 10)",
		})
	}

	// Create post
	post, err := h.service.CreatePost(c.Context(), currentUserID, &req)
	if err != nil {
		if err == ErrInvalidContent {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid content",
			})
		}
		if err == ErrTooManyImages {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "too many images",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetPost handles GET /v1/posts/:postId
func (h *Handler) GetPost(c *fiber.Ctx) error {
	// Get post ID from URL
	postID := c.Params("postId")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post ID is required",
		})
	}

	// Get current user ID (optional, for reaction info)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get post
	post, err := h.service.GetPost(c.Context(), postID, currentUserID)
	if err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get post",
		})
	}

	return c.JSON(post)
}

// GetUserPosts handles GET /v1/users/:userId/posts
func (h *Handler) GetUserPosts(c *fiber.Ctx) error {
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

	// Get posts
	response, err := h.service.GetUserPosts(c.Context(), userID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get posts",
		})
	}

	return c.JSON(response)
}

// DeletePost handles DELETE /v1/posts/:postId
func (h *Handler) DeletePost(c *fiber.Ctx) error {
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

	// Delete post
	err := h.service.DeletePost(c.Context(), postID, currentUserID)
	if err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post not found",
			})
		}
		if err == ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "unauthorized to delete this post",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "post deleted successfully",
	})
}

// ToggleReaction handles POST /v1/posts/:postId/reactions
func (h *Handler) ToggleReaction(c *fiber.Ctx) error {
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

	// Toggle reaction
	response, err := h.service.ToggleReaction(c.Context(), postID, currentUserID)
	if err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to toggle reaction",
		})
	}

	return c.JSON(response)
}

// GetFeed handles GET /v1/feed
func (h *Handler) GetFeed(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Get feed
	response, err := h.service.GetFeed(c.Context(), currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get feed",
		})
	}

	return c.JSON(response)
}
