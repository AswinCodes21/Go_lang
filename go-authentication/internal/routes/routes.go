package routes

import (
	"go-authentication/handlers"
	"go-authentication/internal/constants"
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
		iec104 := auth.Group(constants.IEC104BasePath)
		{
			iec104.GET(constants.IEC104DevicesPath, iec104Handler.GetDevices)
			iec104.GET(constants.IEC104DevicePath+constants.IEC104StatusPath, iec104Handler.GetDeviceStatus)
			iec104.GET(constants.IEC104DevicePath+constants.IEC104DataPath, iec104Handler.GetLatestData)
			iec104.POST(constants.IEC104DevicePath+constants.IEC104ConnectPath, iec104Handler.ConnectToDevice)
			iec104.POST(constants.IEC104DevicePath+constants.IEC104DisconnectPath, iec104Handler.DisconnectFromDevice)
		}
	}
}
