package title

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for title endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new title handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAllTitles handles GET /v1/titles
func (h *Handler) GetAllTitles(c *fiber.Ctx) error {
	// Check for defined_by filter
	definedBy := c.Query("defined_by")

	var titles []Title
	var err error

	if definedBy != "" {
		titles, err = h.service.GetTitlesByDefinedBy(c.Context(), definedBy)
	} else {
		titles, err = h.service.GetAllTitles(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ListTitlesResponse{
		Titles: titles,
		Total:  len(titles),
	})
}

// GetTitle handles GET /v1/titles/:id
func (h *Handler) GetTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title id is required",
		})
	}

	title, err := h.service.GetTitleByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(GetTitleResponse{
		Title: *title,
	})
}

// GetMostUsedTitles handles GET /v1/titles/popular
func (h *Handler) GetMostUsedTitles(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)

	titles, err := h.service.GetMostUsedTitles(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ListTitlesResponse{
		Titles: titles,
		Total:  len(titles),
	})
}

// CreateTitle handles POST /v1/titles
func (h *Handler) CreateTitle(c *fiber.Ctx) error {
	var req CreateTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	title, err := h.service.CreateTitle(c.Context(), req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateTitleResponse{
		Message: "title created successfully",
		Title:   *title,
	})
}

// CreateAdminTitle handles POST /v1/titles/admin (admin only)
func (h *Handler) CreateAdminTitle(c *fiber.Ctx) error {
	var req CreateTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	title, err := h.service.CreateAdminTitle(c.Context(), req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateTitleResponse{
		Message: "admin title created successfully",
		Title:   *title,
	})
}

// UpdateTitle handles PUT /v1/titles/:id
func (h *Handler) UpdateTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title id is required",
		})
	}

	var req UpdateTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	title, err := h.service.UpdateTitle(c.Context(), id, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(UpdateTitleResponse{
		Message: "title updated successfully",
		Title:   *title,
	})
}

// DeleteTitle handles DELETE /v1/titles/:id
func (h *Handler) DeleteTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title id is required",
		})
	}

	if err := h.service.DeleteTitle(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DeleteTitleResponse{
		Message: "title deleted successfully",
	})
}

// IncrementUsage handles POST /v1/titles/:id/increment
func (h *Handler) IncrementUsage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title id is required",
		})
	}

	count, err := h.service.IncrementUsage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(IncrementUsageResponse{
		Message:     "usage incremented successfully",
		UsedByCount: count,
	})
}

// DecrementUsage handles POST /v1/titles/:id/decrement
func (h *Handler) DecrementUsage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title id is required",
		})
	}

	count, err := h.service.DecrementUsage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DecrementUsageResponse{
		Message:     "usage decremented successfully",
		UsedByCount: count,
	})
}
