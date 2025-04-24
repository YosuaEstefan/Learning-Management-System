package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// AssessmentRepository menangani operasi basis data untuk penilaian
type AssessmentRepository struct {
	DB *gorm.DB
}

// NewAssessmentRepository membuat repositori penilaian baru
func NewAssessmentRepository(db *gorm.DB) *AssessmentRepository {
	return &AssessmentRepository{DB: db}
}

// FindByID menemukan penilaian berdasarkan ID
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

// FindBySubmission menemukan penilaian berdasarkan ID pengajuan
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

// Buat membuat penilaian baru
func (r *AssessmentRepository) Create(assessment *models.Assessment) error {
	return r.DB.Create(assessment).Error
}

// Perbarui memperbarui penilaian yang ada
func (r *AssessmentRepository) Update(assessment *models.Assessment) error {
	return r.DB.Save(assessment).Error
}

// Hapus menghapus penilaian
func (r *AssessmentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Assessment{}, id).Error
}
