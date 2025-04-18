package repositories

import (
	"LMS/models"

	"gorm.io/gorm"
)

type AssessmentRepository interface {
	Create(assessment *models.Assessment) error
	GetBySubmission(submissionID uint) (*models.Assessment, error)
}

type assessmentRepository struct {
	db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) AssessmentRepository {
	return &assessmentRepository{db}
}

func (r *assessmentRepository) Create(assessment *models.Assessment) error {
	return r.db.Create(assessment).Error
}

func (r *assessmentRepository) GetBySubmission(submissionID uint) (*models.Assessment, error) {
	var assessment models.Assessment
	err := r.db.Where("submission_id = ?", submissionID).First(&assessment).Error
	return &assessment, err
}
