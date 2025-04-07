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
	"strconv"

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

	// Initialize IEC 104 service
	iec104Port, _ := strconv.Atoi(cfg.IEC104Port)
	timeout, _ := strconv.Atoi(cfg.IEC104Timeout)
	k, _ := strconv.Atoi(cfg.IEC104K)
	w, _ := strconv.Atoi(cfg.IEC104W)

	iec104Service := services.NewIEC104Service(iec104Port, timeout, k, w, natsService)
	if err := iec104Service.Start(); err != nil {
		log.Fatalf("Failed to start IEC 104 service: %v", err)
	}

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
	messageHandler := handlers.NewMessageHandler(natsService, chatUsecase)
	iec104Handler := delivery.NewIEC104Handler(iec104Service)

	// Initialize and configure router
	router := gin.Default()

	// Apply middlewares (if needed, e.g., CORS, logging, recovery)
	// router.Use(someMiddleware())

	// Register routes
	routes.SetupRoutes(router, authHandler, chatHandler, wsHandler, messageHandler, iec104Handler)

	// Start the server
	serverPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	log.Printf("Server running on port %d...", serverPort)
	if err := router.Run(":" + strconv.Itoa(serverPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
