package domain

import (
	"errors"
	// "fmt"
	"regexp"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=10"`
	CreatedAt time.Time `json:"created_at"`
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (u *User) Validate() error {
	// Check for required fields
	if u.Name == "" {
		return errors.New("name field is required and cannot be empty")
	}
	if u.Email == "" {
		return errors.New("email field is required and cannot be empty")
	}
	if u.Password == "" {
		return errors.New("password field is required and cannot be empty")
	}

	// Check password length
	if len(u.Password) < 10 {
		return errors.New("password must be at least 10 characters long")
	}

	// Validate email format
	if !isValidEmail(u.Email) {
		return errors.New("email must be in a valid format (e.g., user@example.com)")
	}

	return nil
}
