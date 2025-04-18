package services

import (
	"LMS/models"
	"LMS/repositories"
)

var AssessmentServiceInstance AssessmentService

type AssessmentService interface {
	CreateAssessment(assessment *models.Assessment) error
	GetAssessmentBySubmission(submissionID uint) (*models.Assessment, error)
}

type assessmentService struct {
	repo repositories.AssessmentRepository
}

func NewAssessmentService(repo repositories.AssessmentRepository) AssessmentService {
	return &assessmentService{repo}
}

func (s *assessmentService) CreateAssessment(assessment *models.Assessment) error {
	return s.repo.Create(assessment)
}

func (s *assessmentService) GetAssessmentBySubmission(submissionID uint) (*models.Assessment, error) {
	return s.repo.GetBySubmission(submissionID)
}
