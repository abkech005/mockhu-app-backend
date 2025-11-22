package auth

import "github.com/gofiber/fiber/v2"

// RegisterRoutes sets up all authentication-related routes.
// It takes the Fiber app and a configured handler with service dependencies.
func RegisterRoutes(app *fiber.App, handler *Handler) {
	auth := app.Group("/v1/auth")

	auth.Post("/signup", handler.Signup)
	auth.Post("/verify", handler.Verify)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
	auth.Post("/resend", handler.Resend)
}
