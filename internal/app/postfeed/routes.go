package postfeed

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers postfeed routes
func RegisterRoutes(app *fiber.App, handler *Handler, engagementHandler *EngagementHandler) {
	postfeeds := app.Group("/v1/postfeeds")

	// Public routes
	postfeeds.Get("/", handler.List)
	postfeeds.Get("/type/:type", handler.GetByType)
	postfeeds.Get("/:id", handler.GetByID)

	// Protected routes
	postfeeds.Post("/", middleware.AuthMiddleware(), handler.Create)
	postfeeds.Put("/:id", middleware.AuthMiddleware(), handler.Update)
	postfeeds.Delete("/:id", middleware.AuthMiddleware(), handler.Delete)

	// Like routes
	postfeeds.Get("/:id/like", engagementHandler.GetLikeStatus)
	postfeeds.Post("/:id/like", middleware.AuthMiddleware(), engagementHandler.LikePostfeed)
	postfeeds.Delete("/:id/like", middleware.AuthMiddleware(), engagementHandler.UnlikePostfeed)

	// Comment routes
	postfeeds.Get("/:id/comments", engagementHandler.ListComments)
	postfeeds.Post("/:id/comments", middleware.AuthMiddleware(), engagementHandler.CreateComment)
	postfeeds.Put("/:id/comments/:comment_id", middleware.AuthMiddleware(), engagementHandler.UpdateComment)
	postfeeds.Delete("/:id/comments/:comment_id", middleware.AuthMiddleware(), engagementHandler.DeleteComment)

	// Share routes
	postfeeds.Get("/:id/shares", engagementHandler.GetShareCount)
	postfeeds.Post("/:id/share", middleware.AuthMiddleware(), engagementHandler.SharePostfeed)

	// User postfeeds route
	users := app.Group("/v1/users")
	users.Get("/:user_id/postfeeds", handler.GetUserPostfeeds)
}
