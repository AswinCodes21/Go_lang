package routes

import (
	"go-authentication/handlers"
	"go-authentication/internal/delivery"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines API routes
func SetupRoutes(router *gin.Engine, authHandler *delivery.AuthHandler, chatHandler *delivery.ChatHandler, wsHandler *delivery.WebSocketHandler, messageHandler *handlers.MessageHandler) {
	// Public routes
	router.POST("/signup", authHandler.SignupHandler)
	router.POST("/login", authHandler.LoginHandler)

	// Protected routes
	auth := router.Group("/")
	auth.Use(delivery.AuthMiddleware())
	{
		// Chat routes
		chat := auth.Group("/chat")
		{
			chat.POST("/send", chatHandler.SendMessageHandler)
			chat.GET("/messages/:user_id", chatHandler.GetConversationMessagesHandler)
		}

		// WebSocket route
		auth.GET("/ws", wsHandler.HandleWebSocket)
	}
}
