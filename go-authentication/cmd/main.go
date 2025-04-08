package main

import (
	"go-authentication/config"
	"go-authentication/db"
	"go-authentication/handlers"
	"go-authentication/internal/delivery"
	"go-authentication/internal/repository"
	"go-authentication/internal/routes"
	"go-authentication/internal/services"
	"go-authentication/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	// This is for local development, in Docker we use environment variables
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	// Load Config
	config.LoadEnv()

	// Initialize config
	cfg := config.LoadConfig()

	// Initialize the DB
	db.ConnectDB(cfg)
	defer db.CloseDB()

	// Run database migrations
	db.Migrate()

	// Initialize NATS connection
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Initialize services
	iec104Service, err := services.NewIEC104Service(natsURL)
	if err != nil {
		log.Fatalf("Failed to initialize IEC104 service: %v", err)
	}

	// Initialize repositories
	userRepository := repository.NewUserRepository()
	chatRepository := repository.NewChatRepository()

	// Initialize usecases
	authUsecase := usecase.NewAuthorizaationcase(userRepository)
	chatUsecase := usecase.NewChatUsecase(chatRepository, userRepository, nil)

	// Initialize handlers
	authHandler := delivery.NewAuthHandler(authUsecase)
	chatHandler := delivery.NewChatHandler(chatUsecase)
	wsHandler := delivery.NewWebSocketHandler(chatUsecase)
	messageHandler := handlers.NewMessageHandler(nil, chatUsecase)
	iec104Handler := delivery.NewIEC104Handler(iec104Service)

	// Initialize and configure router
	router := gin.Default()

	// Apply middlewares (if needed, e.g., CORS, logging, recovery)
	// router.Use(someMiddleware())

	// Register routes
	routes.SetupRoutes(router, authHandler, chatHandler, wsHandler, messageHandler, iec104Handler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server running on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
