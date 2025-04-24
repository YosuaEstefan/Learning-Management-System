package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"

)

// CourseService menangani logika bisnis kursus
type CourseService struct {
	CourseRepo *repositories.CourseRepository
	UserRepo   *repositories.UserRepository
}

// NewCourseService membuat layanan kursus baru
func NewCourseService(courseRepo *repositories.CourseRepository, userRepo *repositories.UserRepository) *CourseService {
	return &CourseService{
		CourseRepo: courseRepo,
		UserRepo:   userRepo,
	}
}

// CreateCourse membuat kursus baru
func (s *CourseService) CreateCourse(course *models.Course) error {
	// Verifikasi keberadaan mentor
	mentor, err := s.UserRepo.FindByID(course.MentorID)
	if err != nil {
		return errors.New("mentor not found")
	}

	// Verifikasi bahwa mentor memiliki peran sebagai mentor
	if mentor.Role != models.RoleMentor && mentor.Role != models.RoleAdmin {
		return errors.New("user is not a mentor")
	}

	// Buat kursus
	return s.CourseRepo.Create(course)
}

// GetCourseByID mendapatkan kursus dengan ID
func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.CourseRepo.FindByIDWithDetails(id)
}

// UpdateCourse memperbarui sebuah kursus
func (s *CourseService) UpdateCourse(course *models.Course) error {
	// Verifikasi keberadaan kursus
	existingCourse, err := s.CourseRepo.FindByID(course.ID)
	if err != nil {
		return err
	}

	// Perbarui hanya bidang yang diizinkan
	existingCourse.Title = course.Title
	existingCourse.Description = course.Description

	return s.CourseRepo.Update(existingCourse)
}

// DeleteCourse menghapus kursus
func (s *CourseService) DeleteCourse(id uint) error {
	return s.CourseRepo.Delete(id)
}

func (s *CourseService) GetAllCourses(limit, offset int) ([]models.Course, error) {
	return s.CourseRepo.ListAll(limit, offset)
}

// GetCoursesByMentor mendapatkan kursus berdasarkan ID mento
func (s *CourseService) GetCoursesByMentor(mentorID uint, limit, offset int) ([]models.Course, error) {
	return s.CourseRepo.ListByMentor(mentorID, limit, offset)
}
