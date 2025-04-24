package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// DiscussionRepository menangani operasi basis data untuk diskusi
type DiscussionRepository struct {
	DB *gorm.DB
}

// NewDiscussionRepository membuat repositori diskusi baru
func NewDiscussionRepository(db *gorm.DB) *DiscussionRepository {
	return &DiscussionRepository{DB: db}
}

// FindByID menemukan diskusi berdasarkan ID
func (r *DiscussionRepository) FindByID(id uint) (*models.Discussion, error) {
	var discussion models.Discussion
	result := r.DB.First(&discussion, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("discussion not found")
		}
		return nil, result.Error
	}
	return &discussion, nil
}

// FindByCourse menemukan diskusi berdasarkan ID kursus
func (r *DiscussionRepository) FindByCourse(courseID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&discussions)
	return discussions, result.Error
}

// FindByUser menemukan diskusi berdasarkan ID pengguna
func (r *DiscussionRepository) FindByUser(userID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.DB.Where("user_id = ?", userID).Preload("Course").Find(&discussions)
	return discussions, result.Error
}

// Membuat membuat diskusi baru
func (r *DiscussionRepository) Create(discussion *models.Discussion) error {
	return r.DB.Create(discussion).Error
}

// Perbarui memperbarui diskusi yang sudah ada
func (r *DiscussionRepository) Update(discussion *models.Discussion) error {
	return r.DB.Save(discussion).Error
}

// Hapus menghapus diskusi
func (r *DiscussionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Discussion{}, id).Error
}
