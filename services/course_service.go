package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
)

// CourseService handles course business logic
type CourseService struct {
	CourseRepo *repositories.CourseRepository
	UserRepo   *repositories.UserRepository
}

// NewCourseService creates a new course service
func NewCourseService(courseRepo *repositories.CourseRepository, userRepo *repositories.UserRepository) *CourseService {
	return &CourseService{
		CourseRepo: courseRepo,
		UserRepo:   userRepo,
	}
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(course *models.Course) error {
	// Verify the mentor exists
	mentor, err := s.UserRepo.FindByID(course.MentorID)
	if err != nil {
		return errors.New("mentor not found")
	}

	// Verify the mentor has the mentor role
	if mentor.Role != models.RoleMentor && mentor.Role != models.RoleAdmin {
		return errors.New("user is not a mentor")
	}

	// Create the course
	return s.CourseRepo.Create(course)
}

// GetCourseByID gets a course by ID
func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.CourseRepo.FindByIDWithDetails(id)
}

// UpdateCourse updates a course
func (s *CourseService) UpdateCourse(course *models.Course) error {
	// Verify the course exists
	existingCourse, err := s.CourseRepo.FindByID(course.ID)
	if err != nil {
		return err
	}

	// Update only allowed fields
	existingCourse.Title = course.Title
	existingCourse.Description = course.Description

	return s.CourseRepo.Update(existingCourse)
}

// DeleteCourse deletes a course
func (s *CourseService) DeleteCourse(id uint) error {
	return s.CourseRepo.Delete(id)
}

// GetAllCourses gets all courses with pagination
func (s *CourseService) GetAllCourses(limit, offset int) ([]models.Course, error) {
	return s.CourseRepo.ListAll(limit, offset)
}

// GetCoursesByMentor gets courses by mentor ID
func (s *CourseService) GetCoursesByMentor(mentorID uint, limit, offset int) ([]models.Course, error) {
	return s.CourseRepo.ListByMentor(mentorID, limit, offset)
}
