package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// MaterialRepository menangani operasi basis data untuk material
type MaterialRepository struct {
	DB *gorm.DB
}

// NewMaterialRepository membuat repositori materi baru
func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{DB: db}
}

// FindByID menemukan materi berdasarkan ID
func (r *MaterialRepository) FindByID(id uint) (*models.Material, error) {
	var material models.Material
	result := r.DB.First(&material, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("material not found")
		}
		return nil, result.Error
	}
	return &material, nil
}

// FindByCourse menemukan materi berdasarkan ID kursus
func (r *MaterialRepository) FindByCourse(courseID uint) ([]models.Material, error) {
	var materials []models.Material
	result := r.DB.Where("course_id = ?", courseID).Find(&materials)
	return materials, result.Error
}

// Membuat membuat materi baru
func (r *MaterialRepository) Create(material *models.Material) error {
	return r.DB.Create(material).Error
}

// Memperbarui materi yang sudah ada
func (r *MaterialRepository) Update(material *models.Material) error {
	return r.DB.Save(material).Error
}

// Hapus menghapus materi
func (r *MaterialRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Material{}, id).Error
}
