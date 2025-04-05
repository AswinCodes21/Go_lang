package delivery

import (
	"context"
	"go-authentication/internal/domain"
	"go-authentication/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles HTTP requests related to chat functionality
type ChatHandler struct {
	ChatUsecase *usecase.ChatUsecase
}

// NewChatHandler creates a new instance of ChatHandler
func NewChatHandler(chatUsecase *usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{ChatUsecase: chatUsecase}
}

// SendMessageHandler handles sending a new message
func (h *ChatHandler) SendMessageHandler(c *gin.Context) {
	// Get sender ID from token
	senderIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert senderID to int (might be float64 from JWT claims)
	var senderID int
	switch v := senderIDValue.(type) {
	case int:
		senderID = v
	case float64:
		senderID = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Parse request body
	var req domain.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Send message
	message, err := h.ChatUsecase.SendMessage(
		context.Background(),
		senderID,
		req.ReceiverID,
		req.Content,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
		"data":    message,
	})
}

// GetConversationMessagesHandler handles retrieving messages between two users
func (h *ChatHandler) GetConversationMessagesHandler(c *gin.Context) {
	// Get user ID from token
	userIDValue, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert userID to int (might be float64 from JWT claims)
	var currentUserID int
	switch v := userIDValue.(type) {
	case int:
		currentUserID = v
	case float64:
		currentUserID = int(v)
	default:
		log.Printf("Invalid user ID type: %T", v)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Get other user's ID from URL parameter
	otherUserIDStr := c.Param("user_id")
	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		log.Printf("Invalid user ID parameter: %s", otherUserIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Ensure we're not trying to get messages with ourselves
	if currentUserID == otherUserID {
		log.Printf("Cannot get messages with self: user_id=%d", currentUserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot get messages with yourself"})
		return
	}

	log.Printf("Getting messages between current user %d and other user %d", currentUserID, otherUserID)

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get messages from database
	messages, err := h.ChatUsecase.GetConversationMessages(
		context.Background(),
		currentUserID,
		otherUserID,
		limit,
		offset,
	)

	if err != nil {
		log.Printf("Error getting messages: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add is_sent flag to each message
	for i := range messages {
		messages[i].IsSent = messages[i].SenderID == currentUserID
	}

	c.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}
