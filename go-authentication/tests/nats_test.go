package tests

import (
	"encoding/json"
	"go-authentication/internal/domain"
	"go-authentication/services"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestNATSCommunication(t *testing.T) {
	// Start single NATS server using Docker
	// docker run -p 4222:4222 nats:latest

	// Create NATS service
	natsService, err := services.NewNatsService()
	if err != nil {
		t.Fatalf("Failed to create NATS service: %v", err)
	}
	defer natsService.Close()

	// Test message
	testMessage := &domain.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello from port 8080!",
	}

	// Subscribe to messages
	received := make(chan *domain.Message)
	err = natsService.SubscribeToPrivateMessages(2, func(msg *domain.Message) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Send message
	err = natsService.SendPrivateMessage(testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Wait for message
	select {
	case msg := <-received:
		if msg.Content != testMessage.Content {
			t.Errorf("Expected message content %s, got %s", testMessage.Content, msg.Content)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestDirectNATSConnection(t *testing.T) {
	// Connect to NATS server
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer nc.Close()

	// Test message
	testMessage := &domain.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello from port 8080!",
	}

	// Subscribe to messages
	received := make(chan *domain.Message)
	_, err = nc.Subscribe("private.2", func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		received <- &message
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish message
	data, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	err = nc.Publish("private.2", data)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for message
	select {
	case msg := <-received:
		if msg.Content != testMessage.Content {
			t.Errorf("Expected message content %s, got %s", testMessage.Content, msg.Content)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}
