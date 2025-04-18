package repositories

import (
	"LMS/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(submission *models.Submission) error
	GetByAssignmentAndStudent(assignmentID, studentID uint) (*models.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}

func (r *submissionRepository) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) GetByAssignmentAndStudent(assignmentID, studentID uint) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).First(&submission).Error
	return &submission, err
}
