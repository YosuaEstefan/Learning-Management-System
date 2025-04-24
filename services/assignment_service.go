package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// AssignmentService handles assignment business logic
type AssignmentService struct {
	AssignmentRepo *repositories.AssignmentRepository
	CourseRepo     *repositories.CourseRepository
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(assignmentRepo *repositories.AssignmentRepository, courseRepo *repositories.CourseRepository) *AssignmentService {
	return &AssignmentService{
		AssignmentRepo: assignmentRepo,
		CourseRepo:     courseRepo,
	}
}

// CreateAssignment creates a new assignment
func (s *AssignmentService) CreateAssignment(assignment *models.Assignment) error {
	// Verify the course exists
	_, err := s.CourseRepo.FindByID(assignment.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Set creation date
	assignment.CreatedAt = time.Now()

	// Create the assignment
	return s.AssignmentRepo.Create(assignment)
}

// GetAssignmentByID gets an assignment by ID
func (s *AssignmentService) GetAssignmentByID(id uint) (*models.Assignment, error) {
	return s.AssignmentRepo.FindByID(id)
}

// GetAssignmentsByCourse gets assignments by course ID
func (s *AssignmentService) GetAssignmentsByCourse(courseID uint) ([]models.Assignment, error) {
	return s.AssignmentRepo.FindByCourse(courseID)
}

// UpdateAssignment updates an assignment
func (s *AssignmentService) UpdateAssignment(assignment *models.Assignment) error {
	// Verify the assignment exists
	existingAssignment, err := s.AssignmentRepo.FindByID(assignment.ID)
	if err != nil {
		return err
	}

	// Update only allowed fields
	existingAssignment.Title = assignment.Title
	existingAssignment.Description = assignment.Description
	existingAssignment.DueDate = assignment.DueDate
	existingAssignment.MaxScore = assignment.MaxScore

	return s.AssignmentRepo.Update(existingAssignment)
}

// DeleteAssignment deletes an assignment
func (s *AssignmentService) DeleteAssignment(id uint) error {
	return s.AssignmentRepo.Delete(id)
}
