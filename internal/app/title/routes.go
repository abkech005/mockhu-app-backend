package title

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all title-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// Title management
	titles := app.Group("/v1/titles")

	// Public routes
	titles.Get("/", handler.GetAllTitles)             // List all titles (with optional ?defined_by=admin filter)
	titles.Get("/popular", handler.GetMostUsedTitles) // Get most used titles
	titles.Get("/:id", handler.GetTitle)              // Get title by ID

	// User routes (create user-defined titles)
	titles.Post("/", handler.CreateTitle) // Create new user-defined title

	// Admin routes (TODO: add admin auth middleware)
	titles.Post("/admin", handler.CreateAdminTitle) // Create admin-defined title
	titles.Put("/:id", handler.UpdateTitle)         // Update title
	titles.Delete("/:id", handler.DeleteTitle)      // Delete title

	// Usage tracking
	titles.Post("/:id/increment", handler.IncrementUsage) // Increment used_by_count
	titles.Post("/:id/decrement", handler.DecrementUsage) // Decrement used_by_count
}
