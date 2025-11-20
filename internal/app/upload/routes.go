package upload

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	handler := NewHandler()

	upload := app.Group("/v1/upload")

	upload.Post("/avatar", handler.Avatar)
}
