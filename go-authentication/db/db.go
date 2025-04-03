package db

import (
	"context"
	"fmt"
	"go-authentication/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ { // Retry mechanism (3 attempts)
		DB, err = pgxpool.New(ctx, dsn)
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Unable to connect to DB: %v", i+1, err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if err != nil {
		log.Fatal("Failed to establish database connection after retries: ", err)
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatal("Database connection test failed: ", err)
	}

	log.Println("Connected to PostgreSQL!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed!")
	}
}
