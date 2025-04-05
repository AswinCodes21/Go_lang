package tests

import (
	"encoding/json"
	"go-authentication/internal/domain"
	"go-authentication/internal/services"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func getNatsURL() string {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		// Check if running in Docker
		if _, err := os.Stat("/.dockerenv"); err == nil {
			natsURL = "nats://nats:4222" // Docker environment
		} else {
			natsURL = "nats://localhost:4222" // Local development
		}
	}
	return natsURL
}

func TestNATSCommunication(t *testing.T) {
	log.Println("Starting NATS Communication Test...")
	log.Println("1. Creating NATS service...")

	// Create NATS service
	natsService, err := services.NewNatsService()
	if err != nil {
		t.Fatalf("Failed to create NATS service: %v", err)
	}
	defer natsService.Close()
	log.Println("✓ NATS service created successfully")

	// Test message
	testMessage := &domain.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello from port 8080!",
	}
	log.Printf("2. Test message created: %+v", testMessage)

	// Subscribe to messages
	log.Println("3. Setting up message subscription...")
	received := make(chan *domain.Message)
	err = natsService.SubscribeToPrivateMessages(2, func(msg *domain.Message) {
		log.Printf("Message received: %+v", msg)
		received <- msg
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	log.Println("✓ Subscription set up successfully")

	// Send message
	log.Println("4. Sending test message...")
	err = natsService.SendPrivateMessage(testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
	log.Println("✓ Message sent successfully")

	// Wait for message
	log.Println("5. Waiting for message receipt...")
	select {
	case msg := <-received:
		log.Printf("Message received: %+v", msg)
		if msg.Content != testMessage.Content {
			t.Errorf("Expected message content %s, got %s", testMessage.Content, msg.Content)
		}
		log.Println("✓ Message content verified")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
	log.Println("✓ NATS Communication Test completed successfully")
}

func TestDirectNATSConnection(t *testing.T) {
	log.Println("Starting Direct NATS Connection Test...")
	log.Println("1. Connecting to NATS server...")

	// Connect to NATS server
	natsURL := getNatsURL()
	log.Printf("Connecting to NATS at: %s", natsURL)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer nc.Close()
	log.Println("✓ Connected to NATS server successfully")

	// Test message
	testMessage := &domain.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello from port 8080!",
	}
	log.Printf("2. Test message created: %+v", testMessage)

	// Subscribe to messages
	log.Println("3. Setting up message subscription...")
	received := make(chan *domain.Message)
	_, err = nc.Subscribe("private.2", func(msg *nats.Msg) {
		var message domain.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}
		log.Printf("Message received: %+v", message)
		received <- &message
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	log.Println("✓ Subscription set up successfully")

	// Publish message
	log.Println("4. Publishing test message...")
	data, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	err = nc.Publish("private.2", data)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	log.Println("✓ Message published successfully")

	// Wait for message
	log.Println("5. Waiting for message receipt...")
	select {
	case msg := <-received:
		log.Printf("Message received: %+v", msg)
		if msg.Content != testMessage.Content {
			t.Errorf("Expected message content %s, got %s", testMessage.Content, msg.Content)
		}
		log.Println("✓ Message content verified")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
	log.Println("✓ Direct NATS Connection Test completed successfully")
}
