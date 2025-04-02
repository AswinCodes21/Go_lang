package services

import (
	"encoding/json"
	"fmt"
	"log"

	"go-authentication/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	nc *nats.Conn // Single NATS server connection
}

func NewNatsService() (*NatsService, error) {
	// Connect to single NATS server
	nc, err := nats.Connect("nats://localhost:4222")
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

func (s *NatsService) SubscribeToPrivateMessages(userID int, callback func(*domain.Message)) error {
	subject := fmt.Sprintf("private.%d", userID)

	_, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}

	return nil
}

func (s *NatsService) SendPrivateMessage(message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	subject := fmt.Sprintf("private.%d", message.ReceiverID)

	if err := s.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (s *NatsService) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
