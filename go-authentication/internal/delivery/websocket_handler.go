package delivery

import (
	"context"
	"encoding/json"
	"go-authentication/internal/domain"
	"go-authentication/internal/usecase"
	"go-authentication/pkg"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time chat
type WebSocketHandler struct {
	ChatUsecase *usecase.ChatUsecase
	// Track active connections
	clients    map[int]*pkg.Client
	clientsMux sync.RWMutex
	// WebSocket upgrader
	upgrader websocket.Upgrader
}

// NewWebSocketHandler creates a new instance of WebSocketHandler
func NewWebSocketHandler(chatUsecase *usecase.ChatUsecase) *WebSocketHandler {
	return &WebSocketHandler{
		ChatUsecase: chatUsecase,
		clients:     make(map[int]*pkg.Client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Allow all origins for development
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocket upgrades the HTTP connection to WebSocket
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from token
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert userID to int (might be float64 from JWT claims)
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create client
	client := pkg.NewClient(conn, userID)

	// Register client
	h.registerClient(client)

	// Subscribe to NATS for real-time messages
	err = h.ChatUsecase.SubscribeToMessages(userID, func(msg *domain.Message) {
		// Send message to client
		msgData, _ := json.Marshal(msg)
		client.Send <- pkg.WebSocketMessage{
			Type: "message",
			Data: msgData,
		}
	})
	if err != nil {
		log.Printf("Error subscribing to messages: %v", err)
	}

	// Start client handlers
	go h.handleMessages(client)
	go h.handleClientConnection(client)
}

// registerClient adds a client to the clients map
func (h *WebSocketHandler) registerClient(client *pkg.Client) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	h.clients[client.ID] = client
	log.Printf("Client connected: %d", client.ID)
}

// unregisterClient removes a client from the clients map
func (h *WebSocketHandler) unregisterClient(client *pkg.Client) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)
		client.Conn.Close()
		log.Printf("Client disconnected: %d", client.ID)
	}
}

// handleMessages handles sending messages to the client
func (h *WebSocketHandler) handleMessages(client *pkg.Client) {
	for message := range client.Send {
		err := client.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to client %d: %v", client.ID, err)
			client.Conn.Close()
			break
		}
	}
}

// handleClientConnection handles reading messages from the client
func (h *WebSocketHandler) handleClientConnection(client *pkg.Client) {
	defer h.unregisterClient(client)

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse message
		var msgReq domain.MessageRequest
		if err := json.Unmarshal(message, &msgReq); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Send message through NATS
		msg := &domain.Message{
			SenderID:   client.ID,
			ReceiverID: msgReq.ReceiverID,
			Content:    msgReq.Content,
		}

		_, err = h.ChatUsecase.SendMessage(context.Background(), msg.SenderID, msg.ReceiverID, msg.Content)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			// Send error message back to client
			errorData, _ := json.Marshal(map[string]string{"error": err.Error()})
			errorMsg := pkg.WebSocketMessage{
				Type: "error",
				Data: errorData,
			}
			client.Send <- errorMsg
		}
	}
}
