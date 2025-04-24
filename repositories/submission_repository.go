package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// SubmissionRepository handles database operations for submissions
type SubmissionRepository struct {
	DB *gorm.DB
}

// NewSubmissionRepository creates a new submission repository
func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{DB: db}
}

// FindByID finds a submission by ID
func (r *SubmissionRepository) FindByID(id uint) (*models.Submission, error) {
	var submission models.Submission
	result := r.DB.First(&submission, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, result.Error
	}
	return &submission, nil
}

// FindByAssignment finds submissions by assignment ID
func (r *SubmissionRepository) FindByAssignment(assignmentID uint) ([]models.Submission, error) {
	var submissions []models.Submission
	result := r.DB.Where("assignment_id = ?", assignmentID).Preload("Student").Find(&submissions)
	return submissions, result.Error
}

// FindByStudent finds submissions by student ID
func (r *SubmissionRepository) FindByStudent(studentID uint) ([]models.Submission, error) {
	var submissions []models.Submission
	result := r.DB.Where("student_id = ?", studentID).Preload("Assignment").Find(&submissions)
	return submissions, result.Error
}

// FindByAssignmentAndStudent finds a submission by assignment ID and student ID
func (r *SubmissionRepository) FindByAssignmentAndStudent(assignmentID, studentID uint) (*models.Submission, error) {
	var submission models.Submission
	result := r.DB.Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).First(&submission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, result.Error
	}
	return &submission, nil
}

// Create creates a new submission
func (r *SubmissionRepository) Create(submission *models.Submission) error {
	return r.DB.Create(submission).Error
}

// Update updates an existing submission
func (r *SubmissionRepository) Update(submission *models.Submission) error {
	return r.DB.Save(submission).Error
}

// Delete deletes a submission
func (r *SubmissionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Submission{}, id).Error
}
