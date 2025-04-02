package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type NatsService struct {
	nc *nats.Conn
}

var natsService *NatsService

func NewNatsService() (*NatsService, error) {
	if natsService != nil {
		return natsService, nil
	}

	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %v", err)
	}

	natsService = &NatsService{
		nc: nc,
	}

	return natsService, nil
}

// GetPrivateSubject creates a unique subject for private messaging between two users
func GetPrivateSubject(user1, user2 string) string {
	// Ensure consistent subject ordering by comparing user IDs
	if user1 < user2 {
		return fmt.Sprintf("chat.private.%s.%s", user1, user2)
	}
	return fmt.Sprintf("chat.private.%s.%s", user2, user1)
}

// SubscribeToPrivateMessages subscribes to private messages for a specific user
func (s *NatsService) SubscribeToPrivateMessages(userID string, messageHandler func(msg *Message)) error {
	// Subscribe to all private messages where this user is either sender or receiver
	pattern := fmt.Sprintf("chat.private.%s.*", userID)
	_, err := s.nc.Subscribe(pattern, func(msg *nats.Msg) {
		var message Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		messageHandler(&message)
	})
	return err
}

// SendPrivateMessage sends a private message from one user to another
func (s *NatsService) SendPrivateMessage(from, to, content string) error {
	message := Message{
		From:    from,
		To:      to,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	subject := GetPrivateSubject(from, to)
	return s.nc.Publish(subject, data)
}

// Close closes the NATS connection
func (s *NatsService) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
