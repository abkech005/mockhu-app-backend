package onboarding

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	handler := NewHandler()

	onboard := app.Group("/v1/onboard")

	onboard.Post("/basic", handler.Basic)
	onboard.Post("/profile", handler.Profile)
	onboard.Post("/interests", handler.Interests)
}
