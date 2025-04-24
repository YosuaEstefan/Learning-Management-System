package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// CourseRepository menangani operasi basis data untuk kursus
type CourseRepository struct {
	DB *gorm.DB
}

// NewCourseRepository membuat repositori kursus baru
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{DB: db}
}

// FindByID menemukan kursus berdasarkan ID
func (r *CourseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course
	result := r.DB.First(&course, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, result.Error
	}
	return &course, nil
}

// FindByIDWithDetails menemukan kursus dengan ID dan memuat data terkait sebelumnya
func (r *CourseRepository) FindByIDWithDetails(id uint) (*models.Course, error) {
	var course models.Course
	result := r.DB.Preload("Mentor").Preload("Materials").Preload("Assignments").First(&course, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, result.Error
	}
	return &course, nil
}

// Membuat membuat kursus baru
func (r *CourseRepository) Create(course *models.Course) error {
	return r.DB.Create(course).Error
}

// Perbarui memperbarui kursus yang sudah ada
func (r *CourseRepository) Update(course *models.Course) error {
	return r.DB.Save(course).Error
}

// Hapus menghapus kursus
func (r *CourseRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Course{}, id).Error
}

// ListAll mencantumkan semua kursus
func (r *CourseRepository) ListAll(limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	result := r.DB.Preload("Mentor").Limit(limit).Offset(offset).Find(&courses)
	return courses, result.Error
}

// ListByMentor mencantumkan kursus berdasarkan ID mentor
func (r *CourseRepository) ListByMentor(mentorID uint, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	result := r.DB.Where("mentor_id = ?", mentorID).Limit(limit).Offset(offset).Find(&courses)
	return courses, result.Error
}
