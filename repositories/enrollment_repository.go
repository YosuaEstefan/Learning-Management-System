package repositories

import (
	"LMS/models"

	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	GetByUserAndCourse(userID, courseID uint) (*models.Enrollment, error)
	GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) GetByUserAndCourse(userID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error
	return &enrollment, err
}

func (r *enrollmentRepository) GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Where("user_id = ?", userID).Find(&enrollments).Error
	return enrollments, err
}
