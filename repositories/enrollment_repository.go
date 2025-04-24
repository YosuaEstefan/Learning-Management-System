package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// EnrollmentRepository menangani operasi basis data untuk pendaftaran
type EnrollmentRepository struct {
	DB *gorm.DB
}

// NewEnrollmentRepository membuat repositori pendaftaran baru
func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{DB: db}
}

// FindByID menemukan pendaftaran dengan ID
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

// FindByUserAndCourse menemukan pendaftaran berdasarkan ID pengguna dan ID kursus
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

// FindByUser menemukan pendaftaran berdasarkan ID pengguna
func (r *EnrollmentRepository) FindByUser(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	result := r.DB.Where("user_id = ?", userID).Preload("Course").Find(&enrollments)
	return enrollments, result.Error
}

// FindByCourse menemukan pendaftaran berdasarkan ID kursus
func (r *EnrollmentRepository) FindByCourse(courseID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	result := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&enrollments)
	return enrollments, result.Error
}

// Buat membuat pendaftaran baru
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.DB.Create(enrollment).Error
}

// Hapus menghapus pendaftaran
func (r *EnrollmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Enrollment{}, id).Error
}
