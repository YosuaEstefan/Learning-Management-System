// services/enrollment_service.go
package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// EnrollmentService handles enrollment business logic
type EnrollmentService struct {
	EnrollmentRepo *repositories.EnrollmentRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(enrollmentRepo *repositories.EnrollmentRepository, userRepo *repositories.UserRepository, courseRepo *repositories.CourseRepository) *EnrollmentService {
	return &EnrollmentService{
		EnrollmentRepo: enrollmentRepo,
		UserRepo:       userRepo,
		CourseRepo:     courseRepo,
	}
}

// CreateEnrollment creates a new enrollment
func (s *EnrollmentService) CreateEnrollment(enrollment *models.Enrollment) error {
	// Verify the user exists
	user, err := s.UserRepo.FindByID(enrollment.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify the user is a student
	if user.Role != models.RoleStudent {
		return errors.New("only students can enroll in courses")
	}

	// Verify the course exists
	_, err = s.CourseRepo.FindByID(enrollment.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Check if the student is already enrolled
	_, err = s.EnrollmentRepo.FindByUserAndCourse(enrollment.UserID, enrollment.CourseID)
	if err == nil {
		return errors.New("student is already enrolled in this course")
	}

	// Set enrollment date
	enrollment.EnrollmentDate = time.Now()

	// Create the enrollment
	return s.EnrollmentRepo.Create(enrollment)
}

// GetEnrollmentByID gets an enrollment by ID
func (s *EnrollmentService) GetEnrollmentByID(id uint) (*models.Enrollment, error) {
	return s.EnrollmentRepo.FindByID(id)
}

// GetEnrollmentsByUser gets enrollments by user ID
func (s *EnrollmentService) GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error) {
	return s.EnrollmentRepo.FindByUser(userID)
}

// GetEnrollmentsByCourse gets enrollments by course ID
func (s *EnrollmentService) GetEnrollmentsByCourse(courseID uint) ([]models.Enrollment, error) {
	return s.EnrollmentRepo.FindByCourse(courseID)
}

// DeleteEnrollment deletes an enrollment
func (s *EnrollmentService) DeleteEnrollment(id uint) error {
	return s.EnrollmentRepo.Delete(id)
}

// IsStudentEnrolled checks if a student is enrolled in a course
func (s *EnrollmentService) IsStudentEnrolled(userID, courseID uint) (bool, error) {
	_, err := s.EnrollmentRepo.FindByUserAndCourse(userID, courseID)
	if err != nil {
		return false, nil
	}
	return true, nil
}
