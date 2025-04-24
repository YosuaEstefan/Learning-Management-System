package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// UserRepository menangani operasi basis data untuk pengguna
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository membuat repositori pengguna baru
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// FindByID menemukan pengguna berdasarkan ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail menemukan pengguna melalui email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Buat membuat pengguna baru
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// Memperbarui memperbarui pengguna yang ada
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// Hapus menghapus pengguna
func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

// ListAll mencantumkan semua pengguna
func (r *UserRepository) ListAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	result := r.DB.Limit(limit).Offset(offset).Find(&users)
	return users, result.Error
}
