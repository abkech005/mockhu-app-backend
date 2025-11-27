package profile

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all profile-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	v1 := app.Group("/v1")

	// Public routes (no auth required)
	v1.Get("/users/:userId/profile", handler.GetUserProfile)
	v1.Get("/users/:userId/mutual-connections", handler.GetMutualConnections)
	v1.Get("/users/:userId/mutual-connections/count", handler.GetMutualConnectionsCount)

	// Protected routes (auth required)
	protected := v1.Group("", middleware.AuthMiddleware())
	protected.Get("/users/me/profile", handler.GetOwnProfile)
	protected.Put("/users/me/profile", handler.UpdateProfile)
	protected.Post("/users/me/avatar", handler.UploadAvatar)
	protected.Delete("/users/me/avatar", handler.DeleteAvatar)
	protected.Get("/users/me/privacy", handler.GetPrivacySettings)
	protected.Put("/users/me/privacy", handler.UpdatePrivacySettings)
}
