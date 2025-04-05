package tests

import (
	"context"
	"errors"
	"go-authentication/internal/domain"
	"go-authentication/internal/services"
	"go-authentication/internal/usecase"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// Mock repositories and services
type mockChatRepo struct {
	messages      []*domain.Message
	conversations []*domain.Conversation
}

func (m *mockChatRepo) StoresMsg(ctx context.Context, message *domain.Message) error {
	if message.Content == "" {
		return errors.New("empty message")
	}
	message.ID = len(m.messages) + 1
	message.CreatedAt = time.Now()
	m.messages = append(m.messages, message)
	return nil
}

func (m *mockChatRepo) GetMsgbyConvo(ctx context.Context, user1ID, user2ID, limit, offset int) ([]*domain.Message, error) {
	if user1ID <= 0 || user2ID <= 0 {
		return nil, errors.New("invalid user IDs")
	}
	var result []*domain.Message
	for _, msg := range m.messages {
		if (msg.SenderID == user1ID && msg.ReceiverID == user2ID) ||
			(msg.SenderID == user2ID && msg.ReceiverID == user1ID) {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (m *mockChatRepo) GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (*domain.Conversation, error) {
	for _, conv := range m.conversations {
		if (conv.User1ID == user1ID && conv.User2ID == user2ID) ||
			(conv.User1ID == user2ID && conv.User2ID == user1ID) {
			return conv, nil
		}
	}
	conv := &domain.Conversation{
		ID:        len(m.conversations) + 1,
		User1ID:   user1ID,
		User2ID:   user2ID,
		UpdatedAt: time.Now(),
	}
	m.conversations = append(m.conversations, conv)
	return conv, nil
}

func (m *mockChatRepo) UpdateConvo(ctx context.Context, conversationID int, lastMessage string) error {
	for _, conv := range m.conversations {
		if conv.ID == conversationID {
			conv.LastMessage = lastMessage
			conv.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("conversation not found")
}

func (m *mockChatRepo) GetConvoByUserId(ctx context.Context, userID int) ([]*domain.Conversation, error) {
	var result []*domain.Conversation
	for _, conv := range m.conversations {
		if conv.User1ID == userID || conv.User2ID == userID {
			result = append(result, conv)
		}
	}
	return result, nil
}

type mockChatUserRepo struct {
	users map[int]*domain.User
}

func (m *mockChatUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockChatUserRepo) Create(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockChatUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// NatsService interface for mocking
type NatsService interface {
	SubscribeToPrivateMessages(userID int, messageHandler func(msg *domain.Message)) error
	SendPrivateMessage(message *domain.Message) error
	SubscribeToSubject(subject string, callback func(*domain.Message)) error
	Close()
}

type mockNatsService struct {
	nc *nats.Conn
}

func NewMockNatsService() *services.NatsService {
	// Create a test-specific NATS service
	ns, err := services.NewNatsService()
	if err != nil {
		log.Printf("Warning: Failed to create NATS service: %v", err)
		// Return a mock service for testing
		return &services.NatsService{}
	}
	return ns
}

func (m *mockNatsService) SubscribeToPrivateMessages(userID int, messageHandler func(msg *domain.Message)) error {
	return nil
}

func (m *mockNatsService) SendPrivateMessage(message *domain.Message) error {
	return nil
}

func (m *mockNatsService) SubscribeToSubject(subject string, callback func(*domain.Message)) error {
	return nil
}

func (m *mockNatsService) Close() {
	// No-op for mock
}

// Test cases
func TestSendMessage(t *testing.T) {
	log.Println("Starting SendMessage Test...")

	// Setup test dependencies
	log.Println("1. Setting up test dependencies...")
	chatRepo := &mockChatRepo{
		messages:      make([]*domain.Message, 0),
		conversations: make([]*domain.Conversation, 0),
	}
	userRepo := &mockChatUserRepo{
		users: map[int]*domain.User{
			1: {ID: 1, Name: "User1"},
			2: {ID: 2, Name: "User2"},
		},
	}
	log.Println("✓ Repositories initialized")

	log.Println("2. Creating NATS service...")
	natsService := NewMockNatsService()
	if natsService == nil {
		t.Fatal("Failed to create NATS service")
	}
	defer natsService.Close()
	log.Println("✓ NATS service created")

	// Create chat usecase with mock dependencies
	log.Println("3. Creating chat usecase...")
	chatUsecase := usecase.NewChatUsecase(chatRepo, userRepo, natsService)
	log.Println("✓ Chat usecase created")

	tests := []struct {
		name       string
		senderID   int
		receiverID int
		content    string
		wantErr    bool
	}{
		{
			name:       "Valid message",
			senderID:   1,
			receiverID: 2,
			content:    "Hello!",
			wantErr:    false,
		},
		{
			name:       "Empty content",
			senderID:   1,
			receiverID: 2,
			content:    "",
			wantErr:    true,
		},
		{
			name:       "Invalid receiver",
			senderID:   1,
			receiverID: 999,
			content:    "Hello!",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("\nRunning test case: %s", tt.name)
			log.Printf("Parameters: sender=%d, receiver=%d, content=%s", tt.senderID, tt.receiverID, tt.content)

			msg, err := chatUsecase.SendMessage(context.Background(), tt.senderID, tt.receiverID, tt.content)
			if (err != nil) != tt.wantErr {
				log.Printf("Test failed: error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && msg == nil {
				log.Println("Test failed: message is nil when no error expected")
				t.Error("SendMessage() returned nil message when no error expected")
			}
			log.Printf("Test passed: %s", tt.name)
		})
	}
	log.Println("✓ SendMessage Test completed")
}

func TestGetMessages(t *testing.T) {
	log.Println("Starting GetMessages Test...")

	// Setup test dependencies with some initial messages
	log.Println("1. Setting up test dependencies...")
	chatRepo := &mockChatRepo{
		messages: []*domain.Message{
			{ID: 1, SenderID: 1, ReceiverID: 2, Content: "Hello"},
			{ID: 2, SenderID: 2, ReceiverID: 1, Content: "Hi"},
		},
	}
	userRepo := &mockChatUserRepo{
		users: map[int]*domain.User{
			1: {ID: 1, Name: "User1"},
			2: {ID: 2, Name: "User2"},
		},
	}
	log.Println("✓ Repositories initialized with test data")

	log.Println("2. Creating NATS service...")
	natsService := NewMockNatsService()
	if natsService == nil {
		t.Fatal("Failed to create NATS service")
	}
	defer natsService.Close()
	log.Println("✓ NATS service created")

	chatUsecase := usecase.NewChatUsecase(chatRepo, userRepo, natsService)

	tests := []struct {
		name      string
		user1ID   int
		user2ID   int
		limit     int
		offset    int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Valid conversation",
			user1ID:   1,
			user2ID:   2,
			limit:     10,
			offset:    0,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Invalid user IDs",
			user1ID:   0,
			user2ID:   2,
			limit:     10,
			offset:    0,
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("\nRunning test case: %s", tt.name)
			log.Printf("Parameters: user1=%d, user2=%d, limit=%d, offset=%d", tt.user1ID, tt.user2ID, tt.limit, tt.offset)

			messages, err := chatUsecase.GetMessages(context.Background(), tt.user1ID, tt.user2ID, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				log.Printf("Test failed: error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(messages) != tt.wantCount {
				log.Printf("Test failed: got %d messages, want %d", len(messages), tt.wantCount)
				t.Errorf("GetMessages() got %d messages, want %d", len(messages), tt.wantCount)
			}
			log.Printf("Test passed: %s", tt.name)
		})
	}
	log.Println("✓ GetMessages Test completed")
}

func TestGetUserConversations(t *testing.T) {
	// Setup test dependencies with some initial conversations
	chatRepo := &mockChatRepo{
		conversations: []*domain.Conversation{
			{ID: 1, User1ID: 1, User2ID: 2, LastMessage: "Hello"},
			{ID: 2, User1ID: 1, User2ID: 3, LastMessage: "Hi"},
		},
	}
	userRepo := &mockChatUserRepo{
		users: map[int]*domain.User{
			1: {ID: 1, Name: "User1"},
			2: {ID: 2, Name: "User2"},
			3: {ID: 3, Name: "User3"},
		},
	}
	natsService := NewMockNatsService()
	if natsService == nil {
		t.Fatal("Failed to create NATS service")
	}
	defer natsService.Close()

	chatUsecase := usecase.NewChatUsecase(chatRepo, userRepo, natsService)

	tests := []struct {
		name      string
		userID    int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "User with conversations",
			userID:    1,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "User without conversations",
			userID:    4,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversations, err := chatUsecase.GetUserConversations(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserConversations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(conversations) != tt.wantCount {
				t.Errorf("GetUserConversations() got %d conversations, want %d", len(conversations), tt.wantCount)
			}
		})
	}
}
