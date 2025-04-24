package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// EnrollmentRepository handles database operations for enrollments
type EnrollmentRepository struct {
	DB *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{DB: db}
}

// FindByID finds an enrollment by ID
func (r *EnrollmentRepository) FindByID(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	result := r.DB.First(&enrollment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, result.Error
	}
	return &enrollment, nil
}

// FindByUserAndCourse finds an enrollment by user ID and course ID
func (r *EnrollmentRepository) FindByUserAndCourse(userID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	result := r.DB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, result.Error
	}
	return &enrollment, nil
}

// FindByUser finds enrollments by user ID
func (r *EnrollmentRepository) FindByUser(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	result := r.DB.Where("user_id = ?", userID).Preload("Course").Find(&enrollments)
	return enrollments, result.Error
}

// FindByCourse finds enrollments by course ID
func (r *EnrollmentRepository) FindByCourse(courseID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	result := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&enrollments)
	return enrollments, result.Error
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.DB.Create(enrollment).Error
}

// Delete deletes an enrollment
func (r *EnrollmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Enrollment{}, id).Error
}
