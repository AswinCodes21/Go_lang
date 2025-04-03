package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"go-authentication/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	nc *nats.Conn
}

func NewNatsService() (*NatsService, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		// Check if running in Docker
		if _, err := os.Stat("/.dockerenv"); err == nil {
			natsURL = "nats://nats:4222" // Docker environment
		} else {
			natsURL = "nats://localhost:4222" // Local development
		}
	}

	// Replace nats:4222 with localhost:4222 if running locally
	if !strings.Contains(natsURL, "localhost") && !strings.Contains(natsURL, "nats:") {
		natsURL = strings.Replace(natsURL, "nats:4222", "localhost:4222", 1)
	}

	log.Printf("Connecting to NATS server at: %s", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS server: %v", err)
	}

	log.Printf("Connected to NATS server successfully!")
	return &NatsService{
		nc: nc,
	}, nil
}

// GetPrivateSubject creates a unique subject for private messaging between two users
func GetPrivateSubject(user1ID, user2ID int) string {
	// Ensure consistent subject ordering by comparing user IDs
	if user1ID < user2ID {
		return fmt.Sprintf("chat.private.%d.%d", user1ID, user2ID)
	}
	return fmt.Sprintf("chat.private.%d.%d", user2ID, user1ID)
}

// SubscribeToPrivateMessages subscribes to private messages for a user
func (s *NatsService) SubscribeToPrivateMessages(userID int, callback func(*domain.Message)) error {
	// Subscribe to all conversations where this user is involved
	subject := fmt.Sprintf("chat.private.*.%d", userID)
	_, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to incoming messages: %w", err)
	}

	// Also subscribe to conversations where user is the second ID
	subject = fmt.Sprintf("chat.private.%d.*", userID)
	_, err = s.nc.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to outgoing messages: %w", err)
	}

	log.Printf("Subscribed to all messages for user %d", userID)
	return nil
}

// SendPrivateMessage sends a message using the unique conversation subject
func (s *NatsService) SendPrivateMessage(message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Get the unique conversation subject
	subject := GetPrivateSubject(message.SenderID, message.ReceiverID)
	log.Printf("Publishing message to subject: %s", subject)

	if err := s.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
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
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}
	return nil
}

func (s *NatsService) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
