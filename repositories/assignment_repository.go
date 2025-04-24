package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// AssignmentRepository handles database operations for assignments
type AssignmentRepository struct {
	DB *gorm.DB
}

// NewAssignmentRepository creates a new assignment repository
func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{DB: db}
}

// FindByID finds an assignment by ID
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

// FindByCourse finds assignments by course ID
func (r *AssignmentRepository) FindByCourse(courseID uint) ([]models.Assignment, error) {
	var assignments []models.Assignment
	result := r.DB.Where("course_id = ?", courseID).Find(&assignments)
	return assignments, result.Error
}

// Create creates a new assignment
func (r *AssignmentRepository) Create(assignment *models.Assignment) error {
	return r.DB.Create(assignment).Error
}

// Update updates an existing assignment
func (r *AssignmentRepository) Update(assignment *models.Assignment) error {
	return r.DB.Save(assignment).Error
}

// Delete deletes an assignment
func (r *AssignmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Assignment{}, id).Error
}
