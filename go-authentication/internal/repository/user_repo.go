package repository

import (
	"context"
	"go-authentication/db"
	"go-authentication/internal/domain"
)

// UserRepository defines the interface for user operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

// userRepository implements UserRepository
type userRepository struct{}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(ctx, query, user.Name, user.Email, user.Password)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	row := db.DB.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE id = $1`
	row := db.DB.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
