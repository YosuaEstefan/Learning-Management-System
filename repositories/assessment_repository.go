package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// AssessmentRepository handles database operations for assessments
type AssessmentRepository struct {
	DB *gorm.DB
}

// NewAssessmentRepository creates a new assessment repository
func NewAssessmentRepository(db *gorm.DB) *AssessmentRepository {
	return &AssessmentRepository{DB: db}
}

// FindByID finds an assessment by ID
func (r *AssessmentRepository) FindByID(id uint) (*models.Assessment, error) {
	var assessment models.Assessment
	result := r.DB.First(&assessment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("assessment not found")
		}
		return nil, result.Error
	}
	return &assessment, nil
}

// FindBySubmission finds an assessment by submission ID
func (r *AssessmentRepository) FindBySubmission(submissionID uint) (*models.Assessment, error) {
	var assessment models.Assessment
	result := r.DB.Where("submission_id = ?", submissionID).First(&assessment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("assessment not found")
		}
		return nil, result.Error
	}
	return &assessment, nil
}

// Create creates a new assessment
func (r *AssessmentRepository) Create(assessment *models.Assessment) error {
	return r.DB.Create(assessment).Error
}

// Update updates an existing assessment
func (r *AssessmentRepository) Update(assessment *models.Assessment) error {
	return r.DB.Save(assessment).Error
}

// Delete deletes an assessment
func (r *AssessmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Assessment{}, id).Error
}
