package services

import (
	"LMS/models"
	"LMS/repositories"
	"LMS/utils"
	"errors"
)

// AuthService handles authentication business logic
type AuthService struct {
	UserRepo *repositories.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

// Register registers a new user
func (s *AuthService) Register(user *models.User) error {
	// Check if user with the same email already exists
	existingUser, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Create the user
	return s.UserRepo.Create(user)
}

// Login logs in a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	// Find the user by email
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Check the password
	if err := user.CheckPassword(password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// GetUserByID gets a user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.UserRepo.FindByID(id)
}
