package handlers

import (
	"net/http"

	"go-authentication/internal/delivery"
	"go-authentication/internal/domain"
	"go-authentication/services"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	natsService *services.NatsService
}

func NewMessageHandler(natsService *services.NatsService) *MessageHandler {
	return &MessageHandler{
		natsService: natsService,
	}
}

type SendMessageRequest struct {
	To      int    `json:"to" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	// Get the sender's ID from the authenticated user context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert float64 to int
	senderID := int(userID.(float64))

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := &domain.Message{
		SenderID:   senderID,
		ReceiverID: req.To,
		Content:    req.Content,
	}

	err := h.natsService.SendPrivateMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func (h *MessageHandler) SetupRoutes(r *gin.Engine) {
	messages := r.Group("/messages")
	messages.Use(delivery.AuthMiddleware())
	{
		messages.POST("/send", h.SendMessage)
	}
}
