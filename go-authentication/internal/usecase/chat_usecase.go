package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-authentication/internal/domain"
	"go-authentication/internal/repository"
	"go-authentication/services"
	"log"
)

// ChatUsecase handles business logic for chat operations
type ChatUsecase struct {
	ChatRepo    repository.ChatRepository
	UserRepo    repository.UserRepository
	NatsService *services.NatsService
}

// NewChatUsecase creates a new instance of ChatUsecase
func NewChatUsecase(chatRepository repository.ChatRepository, userRepository repository.UserRepository, natsService *services.NatsService) *ChatUsecase {
	return &ChatUsecase{
		ChatRepo:    chatRepository,
		UserRepo:    userRepository,
		NatsService: natsService,
	}
}

// SendMessage sends a message from one user to another
func (uc *ChatUsecase) SendMessage(ctx context.Context, senderID int, receiverID int, content string) (*domain.Message, error) {
	log.Printf("Starting SendMessage: sender=%d, receiver=%d, content=%s", senderID, receiverID, content)

	// Validate message content
	if err := domain.ValidateMessage(content); err != nil {
		log.Printf("Invalid message content: %v", err)
		return nil, err
	}

	// Check if receiver exists
	receiver, err := uc.UserRepo.GetByID(ctx, receiverID)
	if err != nil {
		log.Printf("Error checking receiver existence: %v", err)
		return nil, fmt.Errorf("receiver not found: %w", err)
	}

	if receiver == nil {
		log.Printf("Receiver not found with ID: %d", receiverID)
		return nil, errors.New("receiver not found")
	}

	log.Printf("Receiver found: id=%d, name=%s", receiver.ID, receiver.Name)

	// Create message object
	message := &domain.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	log.Printf("Attempting to store message: sender=%d, receiver=%d, content=%s",
		message.SenderID, message.ReceiverID, message.Content)

	// Save message to database
	err = uc.ChatRepo.StoresMsg(ctx, message)
	if err != nil {
		log.Printf("Error storing message: %v", err)
		return nil, fmt.Errorf("failed to store message: %w", err)
	}

	log.Printf("Message stored successfully in database with ID: %d", message.ID)

	// Update conversation
	conversation, err := uc.ChatRepo.GetOrCreateConversation(ctx, senderID, receiverID)
	if err != nil {
		log.Printf("Error getting/creating conversation: %v", err)
		return nil, fmt.Errorf("failed to get/create conversation: %w", err)
	}

	log.Printf("Conversation found/created: id=%d, user1=%d, user2=%d",
		conversation.ID, conversation.User1ID, conversation.User2ID)

	err = uc.ChatRepo.UpdateConvo(ctx, conversation.ID, content)
	if err != nil {
		log.Printf("Error updating conversation: %v", err)
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	log.Printf("Conversation updated with last message")

	// Send message through NATS for real-time delivery
	err = uc.NatsService.SendPrivateMessage(message)
	if err != nil {
		log.Printf("Error sending message through NATS: %v", err)
		// Don't return error here as the message is already saved
	}

	log.Printf("Message sent successfully from user %d to user %d", senderID, receiverID)
	return message, nil
}

// SubscribeToMessages sets up a subscription for real-time messages
func (uc *ChatUsecase) SubscribeToMessages(userID int, callback func(*domain.Message)) error {
	return uc.NatsService.SubscribeToPrivateMessages(userID, callback)
}

// SubscribeToSentMessages subscribes to messages sent BY a user
func (uc *ChatUsecase) SubscribeToSentMessages(userID int, callback func(*domain.Message)) error {
	subject := fmt.Sprintf("chat.sent.%d", userID)
	return uc.NatsService.SubscribeToSubject(subject, callback)
}

// GetConversationMessages retrieves messages between two users
func (uc *ChatUsecase) GetConversationMessages(ctx context.Context, user1ID int, user2ID int, limit int, offset int) ([]*domain.Message, error) {
	log.Printf("Starting GetConversationMessages: user1=%d, user2=%d, limit=%d, offset=%d",
		user1ID, user2ID, limit, offset)

	// Check if user2 exists
	user2, err := uc.UserRepo.GetByID(ctx, user2ID)
	if err != nil || user2 == nil {
		log.Printf("User2 not found: id=%d, error=%v", user2ID, err)
		return nil, errors.New("user not found")
	}

	log.Printf("User2 found: id=%d, name=%s", user2.ID, user2.Name)

	// Set default pagination values
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	log.Printf("Using pagination: limit=%d, offset=%d", limit, offset)

	// Get messages
	messages, err := uc.ChatRepo.GetMsgbyConvo(ctx, user1ID, user2ID, limit, offset)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d messages between users %d and %d", len(messages), user1ID, user2ID)
	return messages, nil
}

// GetUserConversations retrieves all conversations for a user
func (uc *ChatUsecase) GetUserConversations(ctx context.Context, userID int) ([]*domain.Conversation, error) {
	conversations, err := uc.ChatRepo.GetConvoByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

// GetMessages retrieves messages between two users with pagination
func (uc *ChatUsecase) GetMessages(ctx context.Context, user1ID, user2ID int, limit, offset int) ([]*domain.Message, error) {
	log.Printf("Getting messages between users %d and %d (limit: %d, offset: %d)",
		user1ID, user2ID, limit, offset)

	// Get messages from repository
	messages, err := uc.ChatRepo.GetMsgbyConvo(ctx, user1ID, user2ID, limit, offset)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	log.Printf("Found %d messages between users %d and %d", len(messages), user1ID, user2ID)
	return messages, nil
}
