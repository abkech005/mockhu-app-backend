package comment

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for comment operations
type Handler struct {
	service CommentService
}

// NewHandler creates a new comment handler
func NewHandler(service CommentService) *Handler {
	return &Handler{service: service}
}

// CreateComment handles POST /v1/posts/:postId/comments
func (h *Handler) CreateComment(c *fiber.Ctx) error {
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
	var req CreateCommentRequest
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

	if len(req.Content) > 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content too long (max 2000 characters)",
		})
	}

	// Create comment
	comment, err := h.service.CreateComment(c.Context(), postID, currentUserID, &req)
	if err != nil {
		if err == ErrInvalidContent {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid content",
			})
		}
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post not found",
			})
		}
		if err == ErrParentNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "parent comment not found",
			})
		}
		if err == ErrCannotReplyToReply {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot reply to a reply - only top-level comments can have replies",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create comment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

// GetComment handles GET /v1/comments/:commentId
func (h *Handler) GetComment(c *fiber.Ctx) error {
	// Get comment ID from URL
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "comment ID is required",
		})
	}

	// Get current user ID (optional, for author info)
	currentUserID, _ := c.Locals("user_id").(string)

	// Get comment
	comment, err := h.service.GetComment(c.Context(), commentID, currentUserID)
	if err != nil {
		if err == ErrCommentNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "comment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get comment",
		})
	}

	return c.JSON(comment)
}

// GetPostComments handles GET /v1/posts/:postId/comments
func (h *Handler) GetPostComments(c *fiber.Ctx) error {
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

	// Get comments
	response, err := h.service.GetPostComments(c.Context(), postID, currentUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get comments",
		})
	}

	return c.JSON(response)
}

// UpdateComment handles PUT /v1/comments/:commentId
func (h *Handler) UpdateComment(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get comment ID from URL
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "comment ID is required",
		})
	}

	// Parse request body
	var req struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate content
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content is required",
		})
	}

	if len(req.Content) > 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content too long (max 2000 characters)",
		})
	}

	// Update comment
	comment, err := h.service.UpdateComment(c.Context(), commentID, currentUserID, req.Content)
	if err != nil {
		if err == ErrCommentNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "comment not found",
			})
		}
		if err == ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "unauthorized to update this comment",
			})
		}
		if err == ErrInvalidContent {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid content",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update comment",
		})
	}

	return c.JSON(comment)
}

// DeleteComment handles DELETE /v1/comments/:commentId
func (h *Handler) DeleteComment(c *fiber.Ctx) error {
	// Get current user ID from JWT
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get comment ID from URL
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "comment ID is required",
		})
	}

	// Delete comment
	err := h.service.DeleteComment(c.Context(), commentID, currentUserID)
	if err != nil {
		if err == ErrCommentNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "comment not found",
			})
		}
		if err == ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "unauthorized to delete this comment",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete comment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "comment deleted successfully",
	})
}

