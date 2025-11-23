package interest

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all interest-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// Interest management
	interests := app.Group("/v1/interests")
	interests.Get("/", handler.GetAllInterests)         // List all interests (with optional ?category=tech filter)
	interests.Get("/categories", handler.GetCategories) // List all categories with counts
	interests.Post("/", handler.CreateInterest)         // Create new interest (admin)

	// User interests
	users := app.Group("/v1/users")
	users.Get("/:id/interests", handler.GetUserInterests)            // Get user's interests
	users.Post("/:id/interests", handler.AddUserInterests)           // Add interests to user
	users.Put("/:id/interests", handler.ReplaceUserInterests)        // Replace all user interests
	users.Delete("/:id/interests/:slug", handler.RemoveUserInterest) // Remove single interest
}
