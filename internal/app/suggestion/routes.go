package suggestion

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all suggestion-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	suggestions := app.Group("/v1/suggestions")

	// Get user suggestions based on interests
	suggestions.Get("/users", handler.GetUserSuggestions)
}
