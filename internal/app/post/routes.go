package post

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all post-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	v1 := app.Group("/v1")

	// Public routes (no auth required)
	posts := v1.Group("/posts")
	posts.Get("/:postId", handler.GetPost)

	// Protected routes (auth required)
	protected := v1.Group("", middleware.AuthMiddleware())
	protected.Post("/posts", handler.CreatePost)
	protected.Delete("/posts/:postId", handler.DeletePost)
	protected.Post("/posts/:postId/reactions", handler.ToggleReaction)
	protected.Get("/feed", handler.GetFeed)

	// User posts (public, but auth optional for reaction info)
	users := v1.Group("/users")
	users.Get("/:userId/posts", handler.GetUserPosts)
}

