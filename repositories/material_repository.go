package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// MaterialRepository handles database operations for materials
type MaterialRepository struct {
	DB *gorm.DB
}

// NewMaterialRepository creates a new material repository
func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{DB: db}
}

// FindByID finds a material by ID
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

// FindByCourse finds materials by course ID
func (r *MaterialRepository) FindByCourse(courseID uint) ([]models.Material, error) {
	var materials []models.Material
	result := r.DB.Where("course_id = ?", courseID).Find(&materials)
	return materials, result.Error
}

// Create creates a new material
func (r *MaterialRepository) Create(material *models.Material) error {
	return r.DB.Create(material).Error
}

// Update updates an existing material
func (r *MaterialRepository) Update(material *models.Material) error {
	return r.DB.Save(material).Error
}

// Delete deletes a material
func (r *MaterialRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Material{}, id).Error
}
