package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/follow"
	"mockhu-app-backend/internal/app/interest"
	"mockhu-app-backend/internal/app/onboarding"
	"mockhu-app-backend/internal/app/post"
	"mockhu-app-backend/internal/app/upload"
	dbinfra "mockhu-app-backend/internal/infra/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Connect to database
	ctx := context.Background()
	pg, err := dbinfra.New(ctx, dbinfra.DatabaseURLFromEnv())
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer pg.Close()
	log.Println("✅ Database connected")

	app := setupRouter(pg)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		_ = app.Shutdown()
	}()

	log.Println("Server starting on :8085")
	if err := app.Listen(":8085"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func setupRouter(pg *dbinfra.Postgres) *fiber.App {
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

	// Build dependency layers: Repository -> Service -> Handler
	authRepo := auth.NewPostgresUserRepository(pg.Pool)
	verificationRepo := auth.NewPostgresVerificationRepository(pg.Pool)
	authService := auth.NewService(authRepo, verificationRepo)
	authHandler := auth.NewHandler(authService)

	// Interest dependencies
	interestRepo := interest.NewPostgresInterestRepository(pg.Pool)
	interestService := interest.NewService(interestRepo)
	interestHandler := interest.NewHandler(interestService)

	// Onboarding dependencies (reuse authRepo and interestRepo)
	onboardingService := onboarding.NewService(authRepo, interestRepo)
	onboardingHandler := onboarding.NewHandler(onboardingService)

	// Follow dependencies
	followRepo := follow.NewPostgresFollowRepository(pg.Pool)
	followService := follow.NewService(followRepo, authRepo)
	followHandler := follow.NewHandler(followService)

	// Post dependencies
	postRepo := post.NewPostgresPostRepository(pg.Pool)
	postService := post.NewService(postRepo, authRepo)
	postHandler := post.NewHandler(postService)

	// Register domain routes
	auth.RegisterRoutes(app, authHandler)
	interest.RegisterRoutes(app, interestHandler)
	onboarding.RegisterRoutes(app, onboardingHandler)
	upload.RegisterRoutes(app)
	follow.RegisterRoutes(app, followHandler)
	post.RegisterRoutes(app, postHandler)

	return app
}
