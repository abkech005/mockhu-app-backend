package location

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all location-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	locations := app.Group("/v1/locations")

	// Public routes
	locations.Get("/", handler.GetAllLocations)            // List all (optional ?country=India)
	locations.Get("/search", handler.SearchLocations)      // Search/autocomplete ?q=mum
	locations.Get("/popular", handler.GetPopularLocations) // Most used locations
	locations.Get("/:id", handler.GetLocation)             // Get by ID

	// Admin routes
	locations.Post("/", handler.CreateLocation)      // Create location
	locations.Put("/:id", handler.UpdateLocation)    // Update location
	locations.Delete("/:id", handler.DeleteLocation) // Delete location

	// Usage tracking
	locations.Post("/:id/increment", handler.IncrementUsage) // Increment count
	locations.Post("/:id/decrement", handler.DecrementUsage) // Decrement count
}
