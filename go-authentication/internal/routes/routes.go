package routes

import (
	"go-authentication/handlers"
	"go-authentication/internal/delivery"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines API routes
func SetupRoutes(router *gin.Engine, authHandler *delivery.AuthHandler, chatHandler *delivery.ChatHandler, wsHandler *delivery.WebSocketHandler, messageHandler *handlers.MessageHandler, iec104Handler *delivery.IEC104Handler) {
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

		// IEC104 routes
		iec104 := auth.Group("/iec104")
		{
			iec104.GET("/devices", iec104Handler.GetDevices)
			iec104.GET("/device/:device_id/status", iec104Handler.GetDeviceStatus)
			iec104.GET("/device/:device_id/data", iec104Handler.GetLatestData)
			iec104.POST("/device/:device_id/connect", iec104Handler.ConnectToDevice)
			iec104.POST("/device/:device_id/disconnect", iec104Handler.DisconnectFromDevice)
		}
	}
}
