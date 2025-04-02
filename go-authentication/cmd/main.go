package main

import (
	"go-authentication/config"
	"go-authentication/db"
	"go-authentication/handlers"
	"go-authentication/internal/delivery"
	"go-authentication/internal/repository"
	"go-authentication/internal/routes"
	"go-authentication/internal/usecase"
	"go-authentication/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	config.LoadEnv()

	// Initialize config
	cfg := config.LoadConfig()

	// Initialize the DB
	db.ConnectDB(cfg)
	defer db.CloseDB()

	// Run database migrations
	db.Migrate()

	// Initialize NATS service
	natsService, err := services.NewNatsService()
	if err != nil {
		log.Fatalf("Failed to initialize NATS service: %v", err)
	}
	defer natsService.Close()

	// Initialize repositories
	userRepository := repository.NewUserRepository()
	chatRepository := repository.NewChatRepository()

	// Initialize usecases
	authUsecase := usecase.NewAuthorizaationcase(userRepository)
	chatUsecase := usecase.NewChatUsecase(chatRepository, userRepository, natsService)

	// Initialize handlers
	authHandler := delivery.NewAuthHandler(authUsecase)
	chatHandler := delivery.NewChatHandler(chatUsecase)
	wsHandler := delivery.NewWebSocketHandler(chatUsecase)
	messageHandler := handlers.NewMessageHandler(natsService)

	// Initialize and configure router
	router := gin.Default()

	// Apply middlewares (if needed, e.g., CORS, logging, recovery)
	// router.Use(someMiddleware())

	// Register routes
	routes.SetupRoutes(router, authHandler, chatHandler, wsHandler)
	messageHandler.SetupRoutes(router)

	// Start the server
	port := cfg.Port
	log.Printf("Server running on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
