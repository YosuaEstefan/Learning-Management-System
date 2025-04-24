// services/submission_service.go
package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// SubmissionService handles submission business logic
type SubmissionService struct {
	SubmissionRepo *repositories.SubmissionRepository
	AssignmentRepo *repositories.AssignmentRepository
	EnrollmentRepo *repositories.EnrollmentRepository
	UserRepo       *repositories.UserRepository
}

// NewSubmissionService creates a new submission service
func NewSubmissionService(
	submissionRepo *repositories.SubmissionRepository,
	assignmentRepo *repositories.AssignmentRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
	userRepo *repositories.UserRepository,
) *SubmissionService {
	return &SubmissionService{
		SubmissionRepo: submissionRepo,
		AssignmentRepo: assignmentRepo,
		EnrollmentRepo: enrollmentRepo,
		UserRepo:       userRepo,
	}
}

// CreateSubmission creates a new submission
func (s *SubmissionService) CreateSubmission(submission *models.Submission, filePath string) error {
	// Verify the assignment exists
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Verify the student exists
	student, err := s.UserRepo.FindByID(submission.StudentID)
	if err != nil {
		return errors.New("student not found")
	}

	// Verify the student is a student
	if student.Role != models.RoleStudent {
		return errors.New("only students can submit assignments")
	}

	// Verify the student is enrolled in the course
	_, err = s.EnrollmentRepo.FindByUserAndCourse(submission.StudentID, assignment.CourseID)
	if err != nil {
		return errors.New("student is not enrolled in this course")
	}

	// Check for due date if set
	if assignment.DueDate != nil && time.Now().After(*assignment.DueDate) {
		return errors.New("assignment due date has passed")
	}

	// Check if the student has already submitted
	existingSubmission, err := s.SubmissionRepo.FindByAssignmentAndStudent(submission.AssignmentID, submission.StudentID)
	if err == nil && existingSubmission != nil {
		// If submission exists, update it
		existingSubmission.FilePath = filePath
		existingSubmission.SubmittedAt = time.Now()
		return s.SubmissionRepo.Update(existingSubmission)
	}

	// Set the file path and submission date
	submission.FilePath = filePath
	submission.SubmittedAt = time.Now()

	// Create the submission
	return s.SubmissionRepo.Create(submission)
}

// GetSubmissionByID gets a submission by ID
func (s *SubmissionService) GetSubmissionByID(id uint) (*models.Submission, error) {
	return s.SubmissionRepo.FindByID(id)
}

// GetSubmissionsByAssignment gets submissions by assignment ID
func (s *SubmissionService) GetSubmissionsByAssignment(assignmentID uint) ([]models.Submission, error) {
	return s.SubmissionRepo.FindByAssignment(assignmentID)
}

// GetSubmissionsByStudent gets submissions by student ID
func (s *SubmissionService) GetSubmissionsByStudent(studentID uint) ([]models.Submission, error) {
	return s.SubmissionRepo.FindByStudent(studentID)
}

// DeleteSubmission deletes a submission
func (s *SubmissionService) DeleteSubmission(id uint) error {
	return s.SubmissionRepo.Delete(id)
}
