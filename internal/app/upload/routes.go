package upload

import (
	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/infra/r2"
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, r2Client *r2.Client, userRepo auth.UserRepository) {
	service := NewService(r2Client, userRepo)
	handler := NewHandler(service)

	upload := app.Group("/v1/upload")

	// Avatar upload flow
	upload.Post("/avatar/request", handler.RequestAvatarUpload)
	upload.Post("/avatar/confirm", handler.ConfirmAvatarUpload)

	// Media upload for postfeeds (requires auth)
	upload.Post("/media/request", middleware.AuthMiddleware(), handler.RequestMediaUpload)
}
