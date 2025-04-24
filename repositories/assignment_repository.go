package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// AssignmentRepository menangani operasi basis data untuk penugasan
type AssignmentRepository struct {
	DB *gorm.DB
}

// NewAssignmentRepository membuat repositori penugasan baru
func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{DB: db}
}

// FindByID menemukan tugas berdasarkan ID
func (r *AssignmentRepository) FindByID(id uint) (*models.Assignment, error) {
	var assignment models.Assignment
	result := r.DB.First(&assignment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("assignment not found")
		}
		return nil, result.Error
	}
	return &assignment, nil
}

// FindByCourse menemukan tugas berdasarkan ID kursus
func (r *AssignmentRepository) FindByCourse(courseID uint) ([]models.Assignment, error) {
	var assignments []models.Assignment
	result := r.DB.Where("course_id = ?", courseID).Find(&assignments)
	return assignments, result.Error
}

// Membuat membuat tugas baru
func (r *AssignmentRepository) Create(assignment *models.Assignment) error {
	return r.DB.Create(assignment).Error
}

// Perbarui memperbarui tugas yang sudah ada
func (r *AssignmentRepository) Update(assignment *models.Assignment) error {
	return r.DB.Save(assignment).Error
}

// Delete menghapus penugasan
func (r *AssignmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Assignment{}, id).Error
}
