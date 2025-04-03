package tests

import (
	"context"
	"errors"
	"go-authentication/internal/domain"
	"go-authentication/internal/usecase"
	"testing"
	"time"
)

// Mock user repository for testing
type mockAuthUserRepo struct {
	users map[int]*domain.User
}

func (m *mockAuthUserRepo) Create(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockAuthUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockAuthUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func TestSignup(t *testing.T) {
	// Initialize mock repository
	repo := &mockAuthUserRepo{
		users: make(map[int]*domain.User),
	}
	authUsecase := usecase.NewAuthorizaationcase(repo)

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "Valid signup",
			user: &domain.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123456",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			user: &domain.User{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123456",
			},
			wantErr: true,
		},
		{
			name: "Password too short",
			user: &domain.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			user: &domain.User{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authUsecase.Signup(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// Verify user was created
				createdUser, err := repo.GetByEmail(context.Background(), tt.user.Email)
				if err != nil {
					t.Errorf("Failed to get created user: %v", err)
				}
				if createdUser.Name != tt.user.Name || createdUser.Email != tt.user.Email {
					t.Errorf("Created user does not match input: got %+v, want %+v", createdUser, tt.user)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Initialize mock repository with a test user
	repo := &mockAuthUserRepo{
		users: map[int]*domain.User{
			1: {
				ID:        1,
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  "$2a$10$abcdefghijklmnopqrstuvwxyz1234567890", // This is a dummy hash
				CreatedAt: time.Now(),
			},
		},
	}
	authUsecase := usecase.NewAuthorizaationcase(repo)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid login",
			email:    "test@example.com",
			password: "password123456",
			wantErr:  false,
		},
		{
			name:     "Invalid email",
			email:    "nonexistent@example.com",
			password: "password123456",
			wantErr:  true,
		},
		{
			name:     "Empty email",
			email:    "",
			password: "password123456",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			email:    "test@example.com",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authUsecase.Login(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Error("Login() returned empty token when no error expected")
			}
		})
	}
}
