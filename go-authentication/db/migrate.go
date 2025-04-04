package db

import (
	"context"
	"fmt"
	"log"
)

func Migrate() {
	// Drop existing tables if they exist (this will cascade drop all constraints)
	dropTables := `
	DROP TABLE IF EXISTS messages CASCADE;
	DROP TABLE IF EXISTS conversations CASCADE;
	DROP TABLE IF EXISTS users CASCADE;
	`

	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER NOT NULL REFERENCES users(id),
		receiver_id INTEGER NOT NULL REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT different_users CHECK (sender_id != receiver_id)
	);
	`

	conversationsTable := `
	CREATE TABLE IF NOT EXISTS conversations (
		id SERIAL PRIMARY KEY,
		user1_id INTEGER NOT NULL REFERENCES users(id),
		user2_id INTEGER NOT NULL REFERENCES users(id),
		last_message TEXT DEFAULT '',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT conversations_different_users CHECK (user1_id != user2_id),
		CONSTRAINT unique_conversation UNIQUE (user1_id, user2_id)
	);
	`

	// Execute migrations
	migrations := []string{dropTables, usersTable, messagesTable, conversationsTable}

	for _, migration := range migrations {
		_, err := DB.Exec(context.Background(), migration)
		if err != nil {
			log.Fatal("Migration failed: ", err)
		}
	}

	fmt.Println("Database migrations completed successfully!")
}
