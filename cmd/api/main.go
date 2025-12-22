package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/comment"
	"mockhu-app-backend/internal/app/education"
	"mockhu-app-backend/internal/app/follow"
	"mockhu-app-backend/internal/app/interest"
	"mockhu-app-backend/internal/app/messaging"
	"mockhu-app-backend/internal/app/onboarding"
	"mockhu-app-backend/internal/app/post"
	"mockhu-app-backend/internal/app/profile"
	"mockhu-app-backend/internal/app/share"
	"mockhu-app-backend/internal/app/suggestion"
	"mockhu-app-backend/internal/app/title"
	"mockhu-app-backend/internal/app/upload"
	dbinfra "mockhu-app-backend/internal/infra/db"
	"mockhu-app-backend/internal/infra/email"
	"mockhu-app-backend/internal/infra/r2"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or error loading it (using system env vars)")
	}

	// Connect to database
	ctx := context.Background()
	pg, err := dbinfra.New(ctx, dbinfra.DatabaseURLFromEnv())
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer pg.Close()
	log.Println("✅ Database connected")

	// Initialize R2 storage client
	r2Client, err := r2.NewClient(ctx)
	if err != nil {
		log.Printf("⚠️ R2 storage not configured: %v", err)
	} else {
		log.Println("✅ R2 storage connected")
	}

	app := setupRouter(pg, r2Client)

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

func setupRouter(pg *dbinfra.Postgres, r2Client *r2.Client) *fiber.App {
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

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"message":   "Mockhu API is running",
			"timestamp": c.Context().Time().Format("2006-01-02 15:04:05"),
		})
	})

	// Serve static files (avatars)
	app.Static("/avatars", "./storage/avatars")

	// Build dependency layers: Repository -> Service -> Handler
	authRepo := auth.NewPostgresUserRepository(pg.Pool)
	verificationRepo := auth.NewPostgresVerificationRepository(pg.Pool)
	oauthRepo := auth.NewPostgresOAuthRepository(pg.Pool)

	// Initialize OAuth Providers
	providers := make(map[string]auth.OAuthProvider)

	if clientID := os.Getenv("GOOGLE_CLIENT_ID"); clientID != "" {
		providers["google"] = auth.NewGoogleOAuthProvider(
			clientID,
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URL"),
		)
	}

	if clientID := os.Getenv("FACEBOOK_CLIENT_ID"); clientID != "" {
		providers["facebook"] = auth.NewFacebookOAuthProvider(
			clientID,
			os.Getenv("FACEBOOK_CLIENT_SECRET"),
			os.Getenv("FACEBOOK_REDIRECT_URL"),
		)
	}

	if clientID := os.Getenv("APPLE_CLIENT_ID"); clientID != "" {
		providers["apple"] = auth.NewAppleOAuthProvider(
			clientID,
			os.Getenv("APPLE_TEAM_ID"),
			os.Getenv("APPLE_KEY_ID"),
			os.Getenv("APPLE_PRIVATE_KEY"),
			os.Getenv("APPLE_REDIRECT_URL"),
		)
	}

	// Initialize Email Service (SES)
	var emailService auth.EmailService
	sesSender := os.Getenv("SES_SENDER_EMAIL")
	awsRegion := os.Getenv("AWS_REGION")

	if sesSender != "" && awsRegion != "" {
		es, err := email.NewSESEmailService(context.Background(), awsRegion, sesSender)
		if err != nil {
			log.Printf("⚠️ Failed to initialize SES: %v", err)
		} else {
			emailService = es
			log.Println("✅ SES Email Service initialized")
		}
	} else {
		log.Println("⚠️ SES not configured (SES_SENDER_EMAIL or AWS_REGION missing)")
	}

	authService := auth.NewService(authRepo, verificationRepo, oauthRepo, emailService, providers)
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

	// Comment dependencies
	commentRepo := comment.NewPostgresCommentRepository(pg.Pool)
	commentService := comment.NewService(commentRepo, authRepo, postRepo)
	commentHandler := comment.NewHandler(commentService)

	// Share dependencies
	shareRepo := share.NewPostgresShareRepository(pg.Pool)
	shareService := share.NewService(shareRepo, authRepo, postRepo)
	shareHandler := share.NewHandler(shareService)

	// Profile dependencies
	profileRepo := profile.NewPostgresProfileRepository(pg.Pool)
	profileService := profile.NewService(profileRepo, pg.Pool)
	profileHandler := profile.NewHandler(profileService)

	// Messaging dependencies
	convRepo := messaging.NewPostgresConversationRepository(pg.Pool)
	msgRepo := messaging.NewPostgresMessageRepository(pg.Pool)
	blockRepo := messaging.NewPostgresBlockRepository(pg.Pool)
	privacyChecker := messaging.NewPrivacyChecker(authRepo, followRepo, blockRepo)
	messagingService := messaging.NewService(convRepo, msgRepo, blockRepo, authRepo, privacyChecker)
	messagingHandler := messaging.NewHandler(messagingService)

	// Suggestion dependencies
	suggestionRepo := suggestion.NewRepository(pg.Pool)
	suggestionService := suggestion.NewService(suggestionRepo)
	suggestionHandler := suggestion.NewHandler(suggestionService)

	// Title dependencies
	titleRepo := title.NewPostgresTitleRepository(pg.Pool)
	titleService := title.NewService(titleRepo)
	titleHandler := title.NewHandler(titleService)

	// Education dependencies
	educationRepo := education.NewPostgresEducationRepository(pg.Pool)
	educationService := education.NewService(educationRepo)
	educationHandler := education.NewHandler(educationService)

	// Register domain routes
	// Register comment routes BEFORE post routes to avoid route conflicts
	comment.RegisterRoutes(app, commentHandler)
	share.RegisterRoutes(app, shareHandler)
	auth.RegisterRoutes(app, authHandler)
	interest.RegisterRoutes(app, interestHandler)
	onboarding.RegisterRoutes(app, onboardingHandler)
	if r2Client != nil {
		upload.RegisterRoutes(app, r2Client, authRepo)
	}
	follow.RegisterRoutes(app, followHandler)
	post.RegisterRoutes(app, postHandler)
	profile.RegisterRoutes(app, profileHandler)
	messaging.RegisterRoutes(app, messagingHandler)
	suggestion.RegisterRoutes(app, suggestionHandler)
	title.RegisterRoutes(app, titleHandler)
	education.RegisterRoutes(app, educationHandler)

	return app
}

//test
