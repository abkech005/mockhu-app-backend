package onboarding

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all onboarding-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	onboarding := app.Group("/v1/onboarding")

	// Single endpoint for complete onboarding (cost-optimized)
	onboarding.Post("/complete", handler.CompleteOnboarding)

	// Get onboarding status
	onboarding.Get("/status/:user_id", handler.GetOnboardingStatus)
}
