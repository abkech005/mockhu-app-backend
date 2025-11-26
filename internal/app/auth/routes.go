package auth

import "github.com/gofiber/fiber/v2"

// RegisterRoutes sets up all authentication-related routes.
// It takes the Fiber app and a configured handler with service dependencies.
// NOTE: All auth routes are PUBLIC (no AuthMiddleware) - users need to access these without authentication.
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// Create auth group - NO middleware applied (public routes)
	auth := app.Group("/v1/auth")

	// Public auth routes (no authentication required)
	auth.Post("/signup", handler.Signup)    // Public - new user registration
	auth.Post("/login", handler.Login)      // Public - user login
	auth.Post("/verify", handler.Verify)    // Public - email/phone verification
	auth.Post("/refresh", handler.Refresh)  // Public - token refresh
	auth.Post("/logout", handler.Logout)     // Public - logout (stateless)
	auth.Post("/resend", handler.Resend)    // Public - resend verification code

	// Public verification routes (no authentication required)
	auth.Post("/send-email-verification", handler.SendEmailVerification)   // Public
	auth.Post("/verify-email", handler.VerifyEmail)                        // Public
	auth.Post("/send-phone-verification", handler.SendPhoneVerification)   // Public
	auth.Post("/verify-phone", handler.VerifyPhone)                         // Public
}
