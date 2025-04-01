package usecase

import (
	"context"
	"errors"
	"go-authentication/internal/domain"
	"go-authentication/internal/repository"
	"go-authentication/pkg"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	UserRepo repository.UserRepository
}

func NewAuthorizaationcase(userRepository repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{UserRepo: userRepository}
}

// Signup-handler
func (uc *AuthUsecase) Signup(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	existingUser, _ := uc.UserRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return uc.UserRepo.Create(ctx, user)
}

// Login-authenticates a user and returns a JWT token
func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", errors.New("invalid Email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid Email or password")
	}

	token, err := pkg.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}
