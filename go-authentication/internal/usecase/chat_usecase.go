package usecase

import (
	"context"
	"errors"
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
	// Validate message content
	if err := domain.ValidateMessage(content); err != nil {
		return nil, err
	}

	// Check if receiver exists
	receiver, err := uc.UserRepo.GetByID(ctx, receiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}

	if receiver == nil {
		return nil, errors.New("receiver not found")
	}

	// Create message object
	message := &domain.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	// Save message to database
	err = uc.ChatRepo.StoresMsg(ctx, message)
	if err != nil {
		return nil, err
	}

	// Update conversation
	conversation, err := uc.ChatRepo.GetOrCreateConversation(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	err = uc.ChatRepo.UpdateConvo(ctx, conversation.ID, content)
	if err != nil {
		return nil, err
	}

	// Send message through NATS for real-time delivery
	err = uc.NatsService.SendPrivateMessage(message)
	if err != nil {
		// Log the error but don't fail the operation since the message is already saved
		log.Printf("Error sending message through NATS: %v", err)
	}

	return message, nil
}

// SubscribeToMessages sets up a subscription for real-time messages
func (uc *ChatUsecase) SubscribeToMessages(userID int, messageHandler func(msg *domain.Message)) error {
	return uc.NatsService.SubscribeToPrivateMessages(userID, messageHandler)
}

// GetConversationMessages retrieves messages between two users
func (uc *ChatUsecase) GetConversationMessages(ctx context.Context, user1ID int, user2ID int, limit int, offset int) ([]*domain.Message, error) {
	// Check if user2 exists
	user2, err := uc.UserRepo.GetByID(ctx, user2ID)
	if err != nil || user2 == nil {
		return nil, errors.New("user not found")
	}

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

	// Get messages
	messages, err := uc.ChatRepo.GetMsgbyConvo(ctx, user1ID, user2ID, limit, offset)
	if err != nil {
		return nil, err
	}

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
