package routes

import (
	"go-authentication/internal/delivery"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines API routes
func SetupRoutes(router *gin.Engine, authHandler *delivery.AuthHandler, chatHandler *delivery.ChatHandler, wsHandler *delivery.WebSocketHandler) {
	// Public Routes
	router.POST("/signup", authHandler.SignupHandler)
	router.POST("/login", authHandler.LoginHandler)

	// Protected Routes
	router.GET("/protected", delivery.AuthMiddleware(), authHandler.ProtectedHandler)

	// Chat Routes - All require authentication
	chat := router.Group("/chat")
	chat.Use(delivery.AuthMiddleware())
	{
		// Get messages from a conversation with a specific user
		chat.GET("/messages/:user_id", chatHandler.GetConversationMessagesHandler)

		// Get all user's conversations
		chat.GET("/conversations", chatHandler.GetUserConversationsHandler)

		// WebSocket endpoint for real-time chat
		chat.GET("/ws", wsHandler.HandleWebSocket)
	}
}
