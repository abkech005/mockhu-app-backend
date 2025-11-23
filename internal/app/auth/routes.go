package auth

import "github.com/gofiber/fiber/v2"

// RegisterRoutes sets up all authentication-related routes.
// It takes the Fiber app and a configured handler with service dependencies.
func RegisterRoutes(app *fiber.App, handler *Handler) {
	auth := app.Group("/v1/auth")

	// Original auth routes
	auth.Post("/signup", handler.Signup)
	auth.Post("/verify", handler.Verify)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
	auth.Post("/resend", handler.Resend)

	// Verification routes
	auth.Post("/send-email-verification", handler.SendEmailVerification)
	auth.Post("/verify-email", handler.VerifyEmail)
	auth.Post("/send-phone-verification", handler.SendPhoneVerification)
	auth.Post("/verify-phone", handler.VerifyPhone)
}
