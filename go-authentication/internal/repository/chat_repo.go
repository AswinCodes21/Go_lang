package repository

import (
	"context"
	"fmt"
	"go-authentication/db"
	"go-authentication/internal/domain"
	"log"
	"time"
)

// ChatRepository defines the interface for chat-related operations
type ChatRepository interface {
	StoresMsg(ctx context.Context, message *domain.Message) error
	GetMsgbyConvo(ctx context.Context, user1ID, user2ID int, limit, offset int) ([]*domain.Message, error)
	GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (*domain.Conversation, error)
	GetConvoByUserId(ctx context.Context, userID int) ([]*domain.Conversation, error)
	UpdateConvo(ctx context.Context, conversationID int, lastMessage string) error
}

// chatRepository implements ChatRepository
type chatRepository struct{}

// NewChatRepository creates a new instance of chatRepository
func NewChatRepository() ChatRepository {
	return &chatRepository{}
}

// StoresMsg stores a new message in the database
func (r *chatRepository) StoresMsg(ctx context.Context, message *domain.Message) error {
	log.Printf("Storing message: sender_id=%d, receiver_id=%d, content=%s",
		message.SenderID, message.ReceiverID, message.Content)

	// Check if sender and receiver exist
	checkSenderQuery := `SELECT id FROM users WHERE id = $1`
	checkReceiverQuery := `SELECT id FROM users WHERE id = $1`

	var senderID, receiverID int
	err := db.DB.QueryRow(ctx, checkSenderQuery, message.SenderID).Scan(&senderID)
	if err != nil {
		log.Printf("Error checking sender existence: %v", err)
		return fmt.Errorf("sender not found: %w", err)
	}

	err = db.DB.QueryRow(ctx, checkReceiverQuery, message.ReceiverID).Scan(&receiverID)
	if err != nil {
		log.Printf("Error checking receiver existence: %v", err)
		return fmt.Errorf("receiver not found: %w", err)
	}

	log.Printf("Both sender and receiver exist: sender=%d, receiver=%d", senderID, receiverID)

	query := `
		INSERT INTO messages (sender_id, receiver_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err = db.DB.QueryRow(
		ctx,
		query,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		now,
	).Scan(&message.ID)

	if err != nil {
		log.Printf("Error storing message: %v", err)
		return fmt.Errorf("failed to save message: %w", err)
	}

	message.CreatedAt = now
	log.Printf("Message stored successfully with ID: %d, created_at: %v", message.ID, message.CreatedAt)
	return nil
}

// GetMsgbyConvo retrieves messages between two users with pagination
func (r *chatRepository) GetMsgbyConvo(ctx context.Context, user1ID, user2ID int, limit, offset int) ([]*domain.Message, error) {
	log.Printf("Getting messages between users %d and %d (limit: %d, offset: %d)",
		user1ID, user2ID, limit, offset)

	// First, let's check if there are any messages in the table
	checkQuery := `
		SELECT COUNT(*) FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	`
	var count int
	err := db.DB.QueryRow(ctx, checkQuery, user1ID, user2ID).Scan(&count)
	if err != nil {
		log.Printf("Error checking message count: %v", err)
		return nil, fmt.Errorf("failed to check message count: %w", err)
	}
	log.Printf("Total messages found between users: %d", count)

	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	log.Printf("Executing query with parameters: user1=%d, user2=%d, limit=%d, offset=%d",
		user1ID, user2ID, limit, offset)

	rows, err := db.DB.Query(ctx, query, user1ID, user2ID, limit, offset)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		return nil, fmt.Errorf("failed to get messages by conversation: %w", err)
	}
	defer rows.Close()

	messages := []*domain.Message{}
	for rows.Next() {
		msg := &domain.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
		log.Printf("Found message: id=%d, sender=%d, receiver=%d, content=%s, created_at=%v",
			msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over messages: %v", err)
		return nil, fmt.Errorf("error occurred while iterating over messages: %w", err)
	}

	log.Printf("Found %d messages between users %d and %d", len(messages), user1ID, user2ID)
	return messages, nil
}

func (r *chatRepository) GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (*domain.Conversation, error) {
	log.Printf("Getting or creating conversation between users %d and %d", user1ID, user2ID)

	query := `
		SELECT id, user1_id, user2_id, last_message, updated_at
		FROM conversations
		WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)
	`

	conversation := &domain.Conversation{}
	err := db.DB.QueryRow(ctx, query, user1ID, user2ID).Scan(
		&conversation.ID,
		&conversation.User1ID,
		&conversation.User2ID,
		&conversation.LastMessage,
		&conversation.UpdatedAt,
	)

	if err != nil {
		log.Printf("No existing conversation found, creating new one")
		createQuery := `
			INSERT INTO conversations (user1_id, user2_id, last_message, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id, user1_id, user2_id, last_message, updated_at
		`

		now := time.Now()
		err = db.DB.QueryRow(
			ctx,
			createQuery,
			user1ID,
			user2ID,
			"",
			now,
		).Scan(
			&conversation.ID,
			&conversation.User1ID,
			&conversation.User2ID,
			&conversation.LastMessage,
			&conversation.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error creating conversation: %v", err)
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}

		log.Printf("Created new conversation: id=%d, user1=%d, user2=%d",
			conversation.ID, conversation.User1ID, conversation.User2ID)
	} else {
		log.Printf("Found existing conversation: id=%d, user1=%d, user2=%d",
			conversation.ID, conversation.User1ID, conversation.User2ID)
	}

	return conversation, nil
}

func (r *chatRepository) GetConvoByUserId(ctx context.Context, userID int) ([]*domain.Conversation, error) {
	query := `
		SELECT id, user1_id, user2_id, last_message, updated_at
		FROM conversations
		WHERE user1_id = $1 OR user2_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := db.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations by user ID: %w", err)
	}
	defer rows.Close()

	conversations := []*domain.Conversation{}
	for rows.Next() {
		conv := &domain.Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.LastMessage,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over conversations: %w", err)
	}

	return conversations, nil
}

func (r *chatRepository) UpdateConvo(ctx context.Context, conversationID int, lastMessage string) error {
	query := `
		UPDATE conversations
		SET last_message = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := db.DB.Exec(ctx, query, lastMessage, time.Now(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	return nil
}
