package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// SubmissionRepository menangani operasi basis data untuk pengiriman
type SubmissionRepository struct {
	DB *gorm.DB
}

// NewSubmissionRepository membuat repositori pengajuan baru
func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{DB: db}
}

// FindByID menemukan kiriman berdasarkan ID
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

// FindByAssignment menemukan kiriman berdasarkan ID penugasan
func (r *SubmissionRepository) FindByAssignment(assignmentID uint) ([]models.Submission, error) {
	var submissions []models.Submission
	result := r.DB.Where("assignment_id = ?", assignmentID).Preload("Student").Find(&submissions)
	return submissions, result.Error
}

// FindByStudent menemukan kiriman berdasarkan ID siswa
func (r *SubmissionRepository) FindByStudent(studentID uint) ([]models.Submission, error) {
	var submissions []models.Submission
	result := r.DB.Where("student_id = ?", studentID).Preload("Assignment").Find(&submissions)
	return submissions, result.Error
}

// FindByAssignmentAndStudent menemukan kiriman berdasarkan ID tugas dan ID siswa
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

// Buat membuat pengajuan baru
func (r *SubmissionRepository) Create(submission *models.Submission) error {
	return r.DB.Create(submission).Error
}

// Perbarui memperbarui kiriman yang sudah ada
func (r *SubmissionRepository) Update(submission *models.Submission) error {
	return r.DB.Save(submission).Error
}

// Hapus menghapus kiriman
func (r *SubmissionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Submission{}, id).Error
}
