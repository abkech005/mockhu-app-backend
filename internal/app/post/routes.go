package post

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all post-related routes
// This function sets up all endpoints for post operations
// All routes require authentication via JWT middleware
// Parameters:
//   - router: Fiber router group (typically /v1)
//   - handler: Post handler instance with all endpoint implementations
func RegisterRoutes(router fiber.Router, handler *Handler) {
	// Create /posts group
	posts := router.Group("/posts")

	// Protected routes (require authentication)
	// POST /v1/posts - Create a new post
	posts.Post("/", middleware.AuthMiddleware(), handler.CreatePost)

	// GET /v1/posts/:id - Get a single post by ID
	posts.Get("/:id", middleware.AuthMiddleware(), handler.GetPost)

	// DELETE /v1/posts/:id - Delete a post (soft delete)
	posts.Delete("/:id", middleware.AuthMiddleware(), handler.DeletePost)

	// POST /v1/posts/:id/react - Toggle fire reaction on a post
	posts.Post("/:id/react", middleware.AuthMiddleware(), handler.ToggleReaction)

	// GET /v1/posts/user/:userId - Get all posts by a specific user
	posts.Get("/user/:userId", middleware.AuthMiddleware(), handler.GetUserPosts)

	// GET /v1/posts/:id/reactions - Get all reactions for a post
	posts.Get("/:id/reactions", middleware.AuthMiddleware(), handler.GetPostReactions)

	// Feed route (at router level, not under /posts)
	// GET /v1/feed - Get personalized feed for authenticated user
	router.Get("/feed", middleware.AuthMiddleware(), handler.GetFeed)
}
