package services

import (
	"LMS/models"
	"LMS/repositories"
)

var SubmissionServiceInstance SubmissionService

type SubmissionService interface {
	CreateSubmission(submission *models.Submission) error
	GetSubmissionByAssignmentAndStudent(assignmentID, studentID uint) (*models.Submission, error)
}

type submissionService struct {
	repo repositories.SubmissionRepository
}

func NewSubmissionService(repo repositories.SubmissionRepository) SubmissionService {
	return &submissionService{repo}
}

func (s *submissionService) CreateSubmission(submission *models.Submission) error {
	return s.repo.Create(submission)
}

func (s *submissionService) GetSubmissionByAssignmentAndStudent(assignmentID, studentID uint) (*models.Submission, error) {
	return s.repo.GetByAssignmentAndStudent(assignmentID, studentID)
}
