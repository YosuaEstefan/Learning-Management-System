package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// CourseRepository handles database operations for courses
type CourseRepository struct {
	DB *gorm.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{DB: db}
}

// FindByID finds a course by ID
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

// FindByIDWithDetails finds a course by ID and preloads related data
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

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	return r.DB.Create(course).Error
}

// Update updates an existing course
func (r *CourseRepository) Update(course *models.Course) error {
	return r.DB.Save(course).Error
}

// Delete deletes a course
func (r *CourseRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Course{}, id).Error
}

// ListAll lists all courses
func (r *CourseRepository) ListAll(limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	result := r.DB.Preload("Mentor").Limit(limit).Offset(offset).Find(&courses)
	return courses, result.Error
}

// ListByMentor lists courses by mentor ID
func (r *CourseRepository) ListByMentor(mentorID uint, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	result := r.DB.Where("mentor_id = ?", mentorID).Limit(limit).Offset(offset).Find(&courses)
	return courses, result.Error
}
