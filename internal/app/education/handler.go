package education

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for education endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new education handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetUserEducation handles GET /v1/users/:user_id/education
func (h *Handler) GetUserEducation(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	educations, err := h.service.GetUserEducation(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ListEducationResponse{
		Education: educations,
		Total:     len(educations),
	})
}

// GetEducation handles GET /v1/education/:id
func (h *Handler) GetEducation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "education id is required",
		})
	}

	education, err := h.service.GetEducationByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(GetEducationResponse{
		Education: *education,
	})
}

// GetCurrentEducation handles GET /v1/users/:user_id/education/current
func (h *Handler) GetCurrentEducation(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	education, err := h.service.GetCurrentEducation(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if education == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no current education found",
		})
	}

	return c.JSON(GetEducationResponse{
		Education: *education,
	})
}

// CreateEducation handles POST /v1/users/:user_id/education
func (h *Handler) CreateEducation(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	var req CreateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	education, err := h.service.CreateEducation(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateEducationResponse{
		Message:   "education added successfully",
		Education: *education,
	})
}

// UpdateEducation handles PUT /v1/education/:id
func (h *Handler) UpdateEducation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "education id is required",
		})
	}

	var req UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	education, err := h.service.UpdateEducation(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(UpdateEducationResponse{
		Message:   "education updated successfully",
		Education: *education,
	})
}

// DeleteEducation handles DELETE /v1/education/:id
func (h *Handler) DeleteEducation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "education id is required",
		})
	}

	if err := h.service.DeleteEducation(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DeleteEducationResponse{
		Message: "education deleted successfully",
	})
}
