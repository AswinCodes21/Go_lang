package services

import (
	"encoding/json"
	"fmt"
	"log"

	"go-authentication/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	nc1 *nats.Conn
	nc2 *nats.Conn
}

var natsService *NatsService

func NewNatsService() (*NatsService, error) {
	if natsService != nil {
		return natsService, nil
	}

	// Connect to first NATS server
	nc1, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, fmt.Errorf("error connecting to first NATS server: %v", err)
	}

	// Connect to second NATS server
	nc2, err := nats.Connect("nats://localhost:4223")
	if err != nil {
		nc1.Close()
		return nil, fmt.Errorf("error connecting to second NATS server: %v", err)
	}

	natsService = &NatsService{
		nc1: nc1,
		nc2: nc2,
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

	// Subscribe on both servers
	_, err := s.nc1.Subscribe(pattern, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message from server 1: %v", err)
			return
		}
		messageHandler(&message)
	})
	if err != nil {
		return err
	}

	_, err = s.nc2.Subscribe(pattern, func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message from server 2: %v", err)
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

	// Publish to both servers
	if err := s.nc1.Publish(subject, data); err != nil {
		return err
	}
	return s.nc2.Publish(subject, data)
}

// Close closes the NATS connections
func (s *NatsService) Close() {
	if s.nc1 != nil {
		s.nc1.Close()
	}
	if s.nc2 != nil {
		s.nc2.Close()
	}
}
