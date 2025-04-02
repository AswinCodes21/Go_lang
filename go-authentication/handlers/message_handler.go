package handlers

import (
	"net/http"

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
	To      string `json:"to" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	// Get the sender's ID from the authenticated user context
	fromUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.natsService.SendPrivateMessage(fromUser.(string), req.To, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func (h *MessageHandler) SetupRoutes(r *gin.Engine) {
	messages := r.Group("/messages")
	{
		messages.POST("/send", h.SendMessage)
	}
}
