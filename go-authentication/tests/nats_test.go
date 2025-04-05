package tests

import (
	"encoding/json"
	"fmt"
	"go-authentication/internal/domain"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func getNatsURL() string {
	return "nats://localhost:4222" // Always use local NATS server
}

func waitForNatsServer(url string, maxRetries int) error {
	// Check if NATS server is running by trying to connect
	for i := 0; i < maxRetries; i++ {
		nc, err := nats.Connect(url)
		if err == nil {
			nc.Close()
			return nil
		}
		log.Printf("NATS server not ready (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("NATS server not available after %d attempts", maxRetries)
}

func TestNATSCommunication(t *testing.T) {
	log.Println("Starting NATS Communication Test...")

	// 1. Wait for NATS server to be ready
	log.Println("1. Waiting for NATS server to be ready...")
	natsURL := getNatsURL()
	log.Printf("NATS URL: %s", natsURL)

	if err := waitForNatsServer(natsURL, 5); err != nil {
		t.Fatalf("NATS server not ready: %v", err)
	}
	log.Println("✓ NATS server is ready")

	// 2. Connect to NATS server
	log.Println("2. Connecting to NATS server...")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()
	log.Println("✓ Connected to NATS server successfully")

	// 3. Create test message
	msg := &domain.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello from port 8080!",
	}
	log.Printf("3. Test message created: %+v", msg)

	// 4. Set up message subscription
	log.Println("4. Setting up message subscription...")
	received := make(chan *domain.Message, 1) // Buffered channel to prevent blocking
	subject := "chat.private.1.2"             // Using the correct subject format
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
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
	defer sub.Unsubscribe()
	log.Println("✓ Subscription set up successfully")

	// 5. Send message
	log.Println("5. Publishing test message...")
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	err = nc.Publish(subject, data)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	log.Println("✓ Message published successfully")

	// 6. Wait for message
	log.Println("6. Waiting for message receipt...")
	select {
	case receivedMsg := <-received:
		log.Printf("Message received: %+v", receivedMsg)
		if receivedMsg.Content != msg.Content {
			t.Errorf("Expected message content %s, got %s", msg.Content, receivedMsg.Content)
		}
		log.Println("✓ Message content verified")
	case <-time.After(10 * time.Second):
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
	subject := "chat.private.1.2" // Using the correct subject format
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
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

	err = nc.Publish(subject, data)
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
