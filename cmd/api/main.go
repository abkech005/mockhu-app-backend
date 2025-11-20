package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/onboarding"
	"mockhu-app-backend/internal/app/upload"
)

func main() {
	app := setupRouter()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	log.Println("Server starting on :8082")
	if err := app.Listen(":8082"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

func setupRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Mockhu API",
	})

	// Register domain routes
	auth.RegisterRoutes(app)
	onboarding.RegisterRoutes(app)
	upload.RegisterRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return app
}
