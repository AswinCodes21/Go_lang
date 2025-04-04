package delivery

import (
	"context"
	"fmt"
	"go-authentication/internal/domain"
	"go-authentication/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	AuthUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUsecase: authUsecase}
}

// signup-handler
func (h *AuthHandler) SignupHandler(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		// Handle validation errors
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldErr := range validationErr {
				switch fieldErr.Tag() {
				case "required":
					errors[fieldErr.Field()] = fmt.Sprintf("%s is required", fieldErr.Field())
				case "email":
					errors[fieldErr.Field()] = "Please provide a valid email address"
				case "min":
					errors[fieldErr.Field()] = "Password must be at least 10 characters long"
				default:
					errors[fieldErr.Field()] = fmt.Sprintf("Invalid %s", fieldErr.Field())
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	fmt.Printf("Received User: %+v\n", user)

	if err := h.AuthUsecase.Signup(context.Background(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

//login-handler

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})

		return
	}

	token, err := h.AuthUsecase.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *AuthHandler) ProtectedHandler(c *gin.Context) {
	// Get user data from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	emailValue, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found"})
		return
	}

	// Convert userID to int (might be float64 from JWT claims)
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	email, ok := emailValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Protected Data",
		"user_id": userID,
		"email":   email,
	})
}
