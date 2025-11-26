package follow

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all follow-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// Create users group
	users := app.Group("/v1/users")

	// Follow/unfollow (auth required)
	users.Post("/:userId/follow", middleware.AuthMiddleware(), handler.Follow)
	users.Delete("/:userId/follow", middleware.AuthMiddleware(), handler.Unfollow)

	// Get followers/following lists (auth required for is_followed_by_me)
	users.Get("/:userId/followers", middleware.AuthMiddleware(), handler.GetFollowers)
	users.Get("/:userId/following", middleware.AuthMiddleware(), handler.GetFollowing)

	// Check follow status (auth required)
	users.Get("/:userId/is-following", middleware.AuthMiddleware(), handler.IsFollowing)

	// Get follow stats (public, no auth required)
	users.Get("/:userId/follow-stats", handler.GetFollowStats)
}
