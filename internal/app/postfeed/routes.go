package postfeed

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers postfeed routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	postfeeds := app.Group("/v1/postfeeds")

	// Public routes
	postfeeds.Get("/", handler.List)
	postfeeds.Get("/type/:type", handler.GetByType)
	postfeeds.Get("/:id", handler.GetByID)

	// Protected routes
	postfeeds.Post("/", middleware.AuthMiddleware(), handler.Create)
	postfeeds.Put("/:id", middleware.AuthMiddleware(), handler.Update)
	postfeeds.Delete("/:id", middleware.AuthMiddleware(), handler.Delete)

	// User postfeeds route
	users := app.Group("/v1/users")
	users.Get("/:user_id/postfeeds", handler.GetUserPostfeeds)
}
