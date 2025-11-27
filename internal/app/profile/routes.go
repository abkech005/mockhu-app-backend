package profile

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all profile-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	users := app.Group("/v1/users")

	// IMPORTANT: Register literal routes (/me/*) BEFORE parameterized routes (/:userId/*)
	// to avoid route conflicts where :userId matches "me"
	
	// Protected routes (auth required) - literal routes first
	users.Get("/me/profile", middleware.AuthMiddleware(), handler.GetOwnProfile)
	users.Put("/me/profile", middleware.AuthMiddleware(), handler.UpdateProfile)
	users.Post("/me/avatar", middleware.AuthMiddleware(), handler.UploadAvatar)
	users.Delete("/me/avatar", middleware.AuthMiddleware(), handler.DeleteAvatar)
	users.Get("/me/privacy", middleware.AuthMiddleware(), handler.GetPrivacySettings)
	users.Put("/me/privacy", middleware.AuthMiddleware(), handler.UpdatePrivacySettings)

	// Public routes (no auth required) - parameterized routes last
	users.Get("/:userId/profile", handler.GetUserProfile)
	users.Get("/:userId/mutual-connections", handler.GetMutualConnections)
	users.Get("/:userId/mutual-connections/count", handler.GetMutualConnectionsCount)
}
