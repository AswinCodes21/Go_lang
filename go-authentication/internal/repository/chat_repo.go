package repository

import (
	"context"
	"fmt"
	"go-authentication/db"
	"go-authentication/internal/domain"
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
	query := `
		INSERT INTO messages (sender_id, receiver_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err := db.DB.QueryRow(
		ctx,
		query,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		now,
	).Scan(&message.ID)

	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	message.CreatedAt = now
	return nil
}

// GetMsgbyConvo retrieves messages between two users with pagination
func (r *chatRepository) GetMsgbyConvo(ctx context.Context, user1ID, user2ID int, limit, offset int) ([]*domain.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := db.DB.Query(ctx, query, user1ID, user2ID, limit, offset)
	if err != nil {
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
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over messages: %w", err)
	}

	return messages, nil
}

func (r *chatRepository) GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (*domain.Conversation, error) {
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
		createQuery := `
			INSERT INTO conversations (user1_id, user2_id, last_message, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

		now := time.Now()
		err = db.DB.QueryRow(
			ctx,
			createQuery,
			user1ID,
			user2ID,
			"",
			now,
		).Scan(&conversation.ID)

		if err != nil {
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}

		conversation.User1ID = user1ID
		conversation.User2ID = user2ID
		conversation.LastMessage = ""
		conversation.UpdatedAt = now
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
