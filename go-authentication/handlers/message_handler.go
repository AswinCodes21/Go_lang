package handlers

import (
	"net/http"
	"strconv"

	"go-authentication/internal/delivery"
	"go-authentication/internal/usecase"
	"go-authentication/services"

	"log"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	natsService *services.NatsService
	chatUsecase *usecase.ChatUsecase
}

func NewMessageHandler(natsService *services.NatsService, chatUsecase *usecase.ChatUsecase) *MessageHandler {
	return &MessageHandler{
		natsService: natsService,
		chatUsecase: chatUsecase,
	}
}

type SendMessageRequest struct {
	To      int    `json:"to" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// GetMessages retrieves messages between the authenticated user and another user
func (h *MessageHandler) GetMessages(c *gin.Context) {
	// Get the authenticated user's ID
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("No user_id found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert float64 to int
	senderID := int(userID.(float64))
	log.Printf("Authenticated user ID: %d", senderID)

	// Get the other user's ID from URL parameter
	otherUserID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Printf("Invalid user_id parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	log.Printf("Getting messages between users %d and %d (limit: %d, offset: %d)",
		senderID, otherUserID, limit, offset)

	// Get messages through chat usecase
	messages, err := h.chatUsecase.GetMessages(c.Request.Context(), senderID, otherUserID, limit, offset)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	// Set IsSent field for each message
	for _, msg := range messages {
		msg.IsSent = msg.SenderID == senderID
	}

	log.Printf("Found %d messages between users %d and %d", len(messages), senderID, otherUserID)

	// Return the messages
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	// Get the sender's ID from the authenticated user context
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("No user_id found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert float64 to int
	senderID := int(userID.(float64))
	log.Printf("Sender ID from token: %d", senderID)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate receiver ID
	if req.To <= 0 {
		log.Printf("Invalid receiver ID: %d", req.To)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	// Prevent sending message to self
	if senderID == req.To {
		log.Printf("Cannot send message to self: sender=%d, receiver=%d", senderID, req.To)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send message to yourself"})
		return
	}

	log.Printf("Sending message - From: %d, To: %d, Content: %s", senderID, req.To, req.Content)

	// Send message through chat usecase
	message, err := h.chatUsecase.SendMessage(c.Request.Context(), senderID, req.To, req.Content)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	log.Printf("Message sent successfully - ID: %d, Sender: %d, Receiver: %d",
		message.ID, message.SenderID, message.ReceiverID)

	// Return the complete message data
	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
		"data": gin.H{
			"id":          message.ID,
			"sender_id":   message.SenderID,
			"receiver_id": message.ReceiverID,
			"content":     message.Content,
			"created_at":  message.CreatedAt,
		},
	})
}

func (h *MessageHandler) SetupRoutes(r *gin.Engine) {
	messages := r.Group("/messages")
	messages.Use(delivery.AuthMiddleware())
	{
		messages.POST("/send", h.SendMessage)
		messages.GET("/:user_id", h.GetMessages)
	}
}
