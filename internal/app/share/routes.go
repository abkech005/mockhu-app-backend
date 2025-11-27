package share

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all share-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	v1 := app.Group("/v1")

	// Public routes (no auth required) - register first
	v1.Get("/shares/:shareId", handler.GetShare)
	v1.Get("/posts/:postId/shares", handler.GetPostShares)
	v1.Get("/posts/:postId/shares/count", handler.GetShareCount)
	v1.Get("/users/:userId/shares", handler.GetUserShares)

	// Protected routes (auth required)
	protected := v1.Group("/v1/shares", middleware.AuthMiddleware())
	protected.Post("/posts/:postId/shares", handler.CreateShare)
	protected.Delete("/shares/:shareId", handler.DeleteShare)
}
