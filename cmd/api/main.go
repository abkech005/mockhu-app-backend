package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/onboarding"
	"mockhu-app-backend/internal/app/upload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	log.Println("Server starting on :8085")
	if err := app.Listen(":8085"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

func setupRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Mockhu API",
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
	app.Use(recover.New())

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
