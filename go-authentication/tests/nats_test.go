package tests

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func TestNatsConnection() {
	// Connect to first NATS server
	nc1, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error connecting to first NATS server: %v", err)
	}
	defer nc1.Close()

	// Connect to second NATS server
	nc2, err := nats.Connect("nats://localhost:4223")
	if err != nil {
		log.Fatalf("Error connecting to second NATS server: %v", err)
	}
	defer nc2.Close()

	// Subscribe on first server
	subject := "test.subject"
	nc1.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("Received on server 1: %s\n", string(msg.Data))
	})

	// Subscribe on second server
	nc2.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("Received on server 2: %s\n", string(msg.Data))
	})

	// Publish from first server
	message := "Hello from server 1!"
	nc1.Publish(subject, []byte(message))
	fmt.Printf("Published from server 1: %s\n", message)

	// Publish from second server
	message = "Hello from server 2!"
	nc2.Publish(subject, []byte(message))
	fmt.Printf("Published from server 2: %s\n", message)

	// Wait for messages to be received
	time.Sleep(1 * time.Second)
}
