package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go-authentication/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	nc *nats.Conn
}

var natsService *NatsService

func NewNatsService() (*NatsService, error) {
	if natsService != nil {
		return natsService, nil
	}

	// Get NATS URL from environment variable
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://nats:4222" // Default to Docker service name
	}

	// Connect to NATS server with retry logic
	var nc *nats.Conn
	var err error
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		nc, err = nats.Connect(natsURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to NATS (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS after %d attempts: %v", maxRetries, err)
	}

	natsService = &NatsService{
		nc: nc,
	}

	return natsService, nil
}

// GetPrivateSubject creates a unique subject for private messaging between two users
func GetPrivateSubject(user1ID, user2ID int) string {
	// Ensure consistent subject ordering by comparing user IDs
	if user1ID < user2ID {
		return fmt.Sprintf("chat.private.%d.%d", user1ID, user2ID)
	}
	return fmt.Sprintf("chat.private.%d.%d", user2ID, user1ID)
}

// SubscribeToPrivateMessages subscribes to private messages for a specific user
func (s *NatsService) SubscribeToPrivateMessages(userID int, messageHandler func(msg *domain.Message)) error {
	// Subscribe to all private messages where this user is either sender or receiver
	pattern := fmt.Sprintf("chat.private.%d.*", userID)
	_, err := s.nc.Subscribe(pattern, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		messageHandler(&message)
	})
	return err
}

// SendPrivateMessage sends a private message from one user to another
func (s *NatsService) SendPrivateMessage(message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	subject := GetPrivateSubject(message.SenderID, message.ReceiverID)
	return s.nc.Publish(subject, data)
}

// SubscribeToSubject subscribes to a specific NATS subject
func (s *NatsService) SubscribeToSubject(subject string, callback func(*domain.Message)) error {
	_, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	return err
}

// Close closes the NATS connection
func (s *NatsService) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
