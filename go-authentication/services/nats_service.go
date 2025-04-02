package services

import (
	"encoding/json"
	"fmt"
	"log"

	"go-authentication/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	nc1 *nats.Conn // Connection to first NATS server
	nc2 *nats.Conn // Connection to second NATS server
}

func NewNatsService() (*NatsService, error) {
	// Connect to first NATS server
	nc1, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to first NATS server: %v", err)
	}

	// Connect to second NATS server
	nc2, err := nats.Connect("nats://localhost:4223")
	if err != nil {
		nc1.Close()
		return nil, fmt.Errorf("failed to connect to second NATS server: %v", err)
	}

	return &NatsService{
		nc1: nc1,
		nc2: nc2,
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
	// Subscribe to messages from both servers
	subject := fmt.Sprintf("private.%d", userID)

	// Subscribe to first server
	_, err := s.nc1.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to first server: %v", err)
	}

	// Subscribe to second server
	_, err = s.nc2.Subscribe(subject, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		callback(&message)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to second server: %v", err)
	}

	return nil
}

func (s *NatsService) SendPrivateMessage(message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	subject := fmt.Sprintf("private.%d", message.ReceiverID)

	// Publish to both servers
	if err := s.nc1.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish to first server: %v", err)
	}

	if err := s.nc2.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish to second server: %v", err)
	}

	return nil
}

func (s *NatsService) Close() {
	if s.nc1 != nil {
		s.nc1.Close()
	}
	if s.nc2 != nil {
		s.nc2.Close()
	}
}
