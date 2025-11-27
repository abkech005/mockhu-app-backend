package comment

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all comment-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	v1 := app.Group("/v1")

	// Public routes (no auth required) - register first
	v1.Get("/comments/:commentId", handler.GetComment)
	v1.Get("/posts/:postId/comments", handler.GetPostComments)

	// Protected routes (auth required)
	protected := v1.Group("", middleware.AuthMiddleware())
	protected.Post("/posts/:postId/comments", handler.CreateComment)
	protected.Put("/comments/:commentId", handler.UpdateComment)
	protected.Delete("/comments/:commentId", handler.DeleteComment)
}

