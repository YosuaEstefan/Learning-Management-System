package utils

import (
	"LMS/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims mewakili klaim yang disimpan dalam token JWT
// Ini termasuk ID pengguna, email, dan peran pengguna
type JWTClaims struct {
	UserID uint        `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken menghasilkan token JWT baru untuk pengguna
func GenerateToken(user *models.User) (string, error) {
	// dapatkan kunci rahasia dari variabel lingkungan atau gunakan default untuk pengembangan
	secretKey := getEnv("JWT_SECRET_KEY", "your-256-bit-secret")

	// Tetapkan waktu kedaluwarsa
	expirationTime := time.Now().Add(24 * time.Hour)

	// membuat klaim dengan informasi pengguna
	// dan waktu kedaluwarsa
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lms-api",
			Subject:   user.Email,
		},
	}

	// membuat token baru dengan klaim yang ditentukan
	// dan metode penandatanganan HMAC SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// menandatangani token dengan kunci rahasia
	// dan mengembalikan token sebagai string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateToken memvalidasi token JWT dan mengembalikan klaimnya
// Jika token tidak valid, mengembalikan kesalahan
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Dapatkan kunci rahasia dari variabel lingkungan atau gunakan default untuk pengembangan
	secretKey := getEnv("JWT_SECRET_KEY", "your-256-bit-secret")

	// Parse token dengan klaim yang ditentukan
	// dan fungsi untuk memvalidasi metode penandatanganan
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// validasi metode penandatanganan
			// pastikan itu adalah HMAC SHA256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// validasi klaim token
	// jika token valid, kembalikan klaimnya
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// membantu fungsi untuk mendapatkan variabel lingkungan
// dengan nilai default jika tidak ditemukan
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
