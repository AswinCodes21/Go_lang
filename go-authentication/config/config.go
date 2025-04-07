package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiration string
	IEC104Port    string
	IEC104Timeout string
	IEC104K       string
	IEC104W       string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found!")
	}
}

func Getenv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	// Try to load .env file, but don't fail if it's not found
	_ = godotenv.Load()

	return &Config{
		Port:          Getenv("PORT", "8081"),
		DBHost:        Getenv("DB_HOST", "postgres"),
		DBPort:        Getenv("DB_PORT", "5432"),
		DBUser:        Getenv("DB_USER", "postgres"),
		DBPassword:    Getenv("DB_PASSWORD", "admin@123"),
		DBName:        Getenv("DB_NAME", "golang_project"),
		DBSSLMode:     Getenv("DB_SSLMODE", "disable"),
		JWTSecret:     Getenv("JWT_SECRET", "UlVwZFpYbGlzN2N3djd4b2lLMjV6OVF0QzM3TkFqQkY="),
		JWTExpiration: Getenv("JWT_EXPIRATION_HOURS", "24"),
		IEC104Port:    Getenv("IEC104_PORT", "2404"),
		IEC104Timeout: Getenv("IEC104_TIMEOUT", "30"),
		IEC104K:       Getenv("IEC104_K", "12"),
		IEC104W:       Getenv("IEC104_W", "8"),
	}
}
