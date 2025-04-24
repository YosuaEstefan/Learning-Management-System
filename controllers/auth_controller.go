package controllers

import (
	"LMS/models"
	"LMS/services"
	"net/http"

	"github.com/gin-gonic/gin"

)

// AuthController menangani permintaan autentikasi
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController membuat pengontrol autentikasi baru
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// RegisterRequest mewakili permintaan untuk mendaftarkan pengguna baru
type RegisterRequest struct {
	Name     string      `json:"name" binding:"required"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=6"`
	Role     models.Role `json:"role" binding:"required,oneof=admin mentor student"`
}

// LoginRequest mewakili permintaan untuk login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (c *AuthController) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	}

	if err := c.AuthService.Register(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login menangani login pengguna
func (c *AuthController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.AuthService.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// GetProfile mengembalikan profil pengguna saat ini
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	user, err := c.AuthService.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
