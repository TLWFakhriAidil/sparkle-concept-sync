package main

import (
	"log"
	"os"

	"sparkle-concept-sync/internal/config"
	"sparkle-concept-sync/internal/database"
	"sparkle-concept-sync/internal/handlers"
	"sparkle-concept-sync/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis service
	redisService := services.NewRedisService(cfg.RedisURL)

	// Initialize services
	aiService := services.NewAIService(cfg.OpenRouterAPIKey, redisService)
	deviceService := services.NewDeviceSettingsService(db)
	flowService := services.NewFlowService(db, aiService)
	providerService := services.NewProviderService()
	websocketService := services.NewWebSocketService()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:    50 * 1024 * 1024, // 50MB limit for media uploads
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Serve static files
	app.Static("/", "./dist")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	profileHandler := handlers.NewProfileHandler(db)
	deviceHandler := handlers.NewDeviceSettingsHandler(deviceService)
	healthHandler := handlers.NewHealthHandler(db, redisService)
	wahaHandler := handlers.NewWAHAHandler(flowService, providerService, websocketService)

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/ws", websocket.New(websocketService.HandleWebSocket))

	// Health check routes
	app.Get("/healthz", healthHandler.HealthCheck)
	app.Get("/api/health/detailed", healthHandler.DetailedHealthCheck)

	// API routes
	api := app.Group("/api")

	// Authentication routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes middleware
	api.Use(authHandler.JWTMiddleware())

	// Profile routes
	profile := api.Group("/profile")
	profile.Get("/", profileHandler.GetProfile)
	profile.Put("/", profileHandler.UpdateProfile)

	// Device settings routes
	devices := api.Group("/device-settings")
	devices.Get("/", deviceHandler.GetDevices)
	devices.Post("/", deviceHandler.CreateDevice)
	devices.Get("/:id", deviceHandler.GetDevice)
	devices.Put("/:id", deviceHandler.UpdateDevice)
	devices.Delete("/:id", deviceHandler.DeleteDevice)

	// Webhook info route
	api.Get("/webhook-info", wahaHandler.GetWebhookInfo)

	// Webhook routes (no auth required)
	webhooks := app.Group("/webhooks")
	webhooks.Post("/waha/:device_id", wahaHandler.HandleWAHAWebhook)
	webhooks.Post("/wablas/:device_id", wahaHandler.HandleWablasWebhook)
	webhooks.Post("/whacenter/:device_id", wahaHandler.HandleWhacenterWebhook)

	// TODO: Media upload routes to be implemented
	// media := api.Group("/media")
	// media.Post("/upload", mediaHandler.HandleUpload)

	// Fallback to serve React app
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./dist/index.html")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Max concurrent users: %d", cfg.MaxConcurrentUsers)
	log.Printf("🔄 WebSocket enabled: %v", cfg.WebSocketEnabled)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
