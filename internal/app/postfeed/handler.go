package postfeed

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for postfeed endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new postfeed handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /v1/postfeeds
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req CreatePostfeedRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	postfeed, err := h.service.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreatePostfeedResponse{
		Message:  "postfeed created successfully",
		Postfeed: *postfeed,
	})
}

// GetByID handles GET /v1/postfeeds/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "postfeed id is required",
		})
	}

	postfeed, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(postfeed)
}

// Update handles PUT /v1/postfeeds/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "postfeed id is required",
		})
	}

	var req UpdatePostfeedRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	postfeed, err := h.service.Update(c.Context(), id, userID, req)
	if err != nil {
		if err.Error() == "unauthorized: you can only update your own posts" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(UpdatePostfeedResponse{
		Message:  "postfeed updated successfully",
		Postfeed: *postfeed,
	})
}

// Delete handles DELETE /v1/postfeeds/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "postfeed id is required",
		})
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		if err.Error() == "unauthorized: you can only delete your own posts" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DeletePostfeedResponse{
		Message: "postfeed deleted successfully",
	})
}

// List handles GET /v1/postfeeds
func (h *Handler) List(c *fiber.Ctx) error {
	req := ListPostfeedsRequest{
		Type:   c.Query("type"),
		UserID: c.Query("user_id"),
		Tag:    c.Query("tag"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
	}

	result, err := h.service.List(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetByType handles GET /v1/postfeeds/type/:type
func (h *Handler) GetByType(c *fiber.Ctx) error {
	postType := c.Params("type")
	if postType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type is required",
		})
	}

	if !IsValidType(postType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "invalid post type",
			"valid_types": ValidTypes(),
		})
	}

	req := ListPostfeedsRequest{
		Type:  postType,
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 20),
	}

	result, err := h.service.List(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetUserPostfeeds handles GET /v1/users/:user_id/postfeeds
func (h *Handler) GetUserPostfeeds(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	result, err := h.service.GetUserPostfeeds(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
