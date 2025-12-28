package education

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all education-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// User education routes
	users := app.Group("/v1/users")
	users.Get("/:user_id/education", handler.GetUserEducation)            // Get all education for a user
	users.Get("/:user_id/education/current", handler.GetCurrentEducation) // Get current education
	users.Post("/:user_id/education", handler.CreateEducation)            // Add education entry

	// Education entry routes (by education ID)
	education := app.Group("/v1/education")
	education.Get("/:id", handler.GetEducation)       // Get education by ID
	education.Put("/:id", handler.UpdateEducation)    // Update education entry
	education.Delete("/:id", handler.DeleteEducation) // Delete education entry
}
