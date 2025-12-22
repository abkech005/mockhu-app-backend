package suggestion

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for suggestions
type Handler struct {
	service *Service
}

// NewHandler creates a new suggestion handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetUserSuggestions returns suggested users based on interests
// GET /v1/suggestions/users?user_id=xxx&limit=10
func (h *Handler) GetUserSuggestions(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id query parameter is required",
		})
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	response, err := h.service.GetUserSuggestions(c.Context(), userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
