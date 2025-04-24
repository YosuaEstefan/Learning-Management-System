package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// AssessmentService handles assessment business logic
type AssessmentService struct {
	AssessmentRepo *repositories.AssessmentRepository
	SubmissionRepo *repositories.SubmissionRepository
	AssignmentRepo *repositories.AssignmentRepository
	UserRepo       *repositories.UserRepository
}

// NewAssessmentService creates a new assessment service
func NewAssessmentService(
	assessmentRepo *repositories.AssessmentRepository,
	submissionRepo *repositories.SubmissionRepository,
	assignmentRepo *repositories.AssignmentRepository,
	userRepo *repositories.UserRepository,
) *AssessmentService {
	return &AssessmentService{
		AssessmentRepo: assessmentRepo,
		SubmissionRepo: submissionRepo,
		AssignmentRepo: assignmentRepo,
		UserRepo:       userRepo,
	}
}

// CreateAssessment creates a new assessment
func (s *AssessmentService) CreateAssessment(assessment *models.Assessment) error {
	// Verify the submission exists
	submission, err := s.SubmissionRepo.FindByID(assessment.SubmissionID)
	if err != nil {
		return errors.New("submission not found")
	}

	// Verify the assignment exists
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Validate the score if max score is set
	if assignment.MaxScore != nil && assessment.Score != nil {
		if *assessment.Score < 0 || *assessment.Score > *assignment.MaxScore {
			return errors.New("score must be between 0 and the maximum score")
		}
	}

	// Check if assessment already exists
	existingAssessment, err := s.AssessmentRepo.FindBySubmission(assessment.SubmissionID)
	if err == nil && existingAssessment != nil {
		// If assessment exists, update it
		existingAssessment.Score = assessment.Score
		existingAssessment.Feedback = assessment.Feedback
		existingAssessment.AssessedAt = time.Now()
		return s.AssessmentRepo.Update(existingAssessment)
	}

	// Set assessment date
	assessment.AssessedAt = time.Now()

	// Create the assessment
	return s.AssessmentRepo.Create(assessment)
}

// GetAssessmentByID gets an assessment by ID
func (s *AssessmentService) GetAssessmentByID(id uint) (*models.Assessment, error) {
	return s.AssessmentRepo.FindByID(id)
}

// GetAssessmentBySubmission gets an assessment by submission ID
func (s *AssessmentService) GetAssessmentBySubmission(submissionID uint) (*models.Assessment, error) {
	return s.AssessmentRepo.FindBySubmission(submissionID)
}

// UpdateAssessment updates an assessment
func (s *AssessmentService) UpdateAssessment(assessment *models.Assessment) error {
	// Verify the assessment exists
	existingAssessment, err := s.AssessmentRepo.FindByID(assessment.ID)
	if err != nil {
		return err
	}

	// Verify the submission exists
	submission, err := s.SubmissionRepo.FindByID(existingAssessment.SubmissionID)
	if err != nil {
		return errors.New("submission not found")
	}

	// Verify the assignment exists
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Validate the score if max score is set
	if assignment.MaxScore != nil && assessment.Score != nil {
		if *assessment.Score < 0 || *assessment.Score > *assignment.MaxScore {
			return errors.New("score must be between 0 and the maximum score")
		}
	}

	// Update only allowed fields
	existingAssessment.Score = assessment.Score
	existingAssessment.Feedback = assessment.Feedback
	existingAssessment.AssessedAt = time.Now()

	return s.AssessmentRepo.Update(existingAssessment)
}

// DeleteAssessment deletes an assessment
func (s *AssessmentService) DeleteAssessment(id uint) error {
	return s.AssessmentRepo.Delete(id)
}
