package services

import (
	"LMS/models"
	"LMS/repositories"
	"LMS/utils"
	"errors"

)

// AuthService menangani logika bisnis otentikasi
type AuthService struct {
	UserRepo *repositories.UserRepository
}

// NewAuthService membuat layanan otentikasi baru
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

// Daftar mendaftarkan pengguna baru
func (s *AuthService) Register(user *models.User) error {
	// Periksa apakah pengguna dengan email yang sama sudah ada
	existingUser, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}
	// Buat pengguna
	return s.UserRepo.Create(user)
}

// Login memasukkan pengguna dan mengembalikan token JWT
func (s *AuthService) Login(email, password string) (string, error) {
	// Temukan pengguna melalui email
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Periksa kata sandi
	if err := user.CheckPassword(password); err != nil {
		return "", errors.New("invalid email or password")
	}
	// Menghasilkan token JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// GetUserByID mendapatkan pengguna dengan ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.UserRepo.FindByID(id)
}
