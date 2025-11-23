package onboarding

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for onboarding endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new onboarding handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CompleteOnboarding handles the single onboarding request
// POST /v1/onboarding/complete
func (h *Handler) CompleteOnboarding(c *fiber.Ctx) error {
	var req CompleteOnboardingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	response, err := h.service.CompleteOnboarding(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetOnboardingStatus returns the current onboarding status
// GET /v1/onboarding/status/:user_id
func (h *Handler) GetOnboardingStatus(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	response, err := h.service.GetOnboardingStatus(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
