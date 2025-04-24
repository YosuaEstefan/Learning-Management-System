package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// AssessmentService menangani logika bisnis penilaian
type AssessmentService struct {
	AssessmentRepo *repositories.AssessmentRepository
	SubmissionRepo *repositories.SubmissionRepository
	AssignmentRepo *repositories.AssignmentRepository
	UserRepo       *repositories.UserRepository
}

// NewAssessmentService membuat layanan penilaian baru
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

// CreateAssessment membuat penilaian baru
func (s *AssessmentService) CreateAssessment(assessment *models.Assessment) error {
	// Verifikasi adanya kiriman
	submission, err := s.SubmissionRepo.FindByID(assessment.SubmissionID)
	if err != nil {
		return errors.New("submission not found")
	}

	// Verifikasi bahwa penugasan itu ada
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Validasi skor jika skor maksimal ditetapkan
	if assignment.MaxScore != nil && assessment.Score != nil {
		if *assessment.Score < 0 || *assessment.Score > *assignment.MaxScore {
			return errors.New("score must be between 0 and the maximum score")
		}
	}

	// Periksa apakah penilaian sudah ada
	existingAssessment, err := s.AssessmentRepo.FindBySubmission(assessment.SubmissionID)
	if err == nil && existingAssessment != nil {
		// Jika penilaian ada, perbarui
		existingAssessment.Score = assessment.Score
		existingAssessment.Feedback = assessment.Feedback
		existingAssessment.AssessedAt = time.Now()
		return s.AssessmentRepo.Update(existingAssessment)
	}

	// Tetapkan tanggal penilaian
	assessment.AssessedAt = time.Now()

	// Buat penilaian
	return s.AssessmentRepo.Create(assessment)
}

// GetAssessmentByID mendapatkan penilaian berdasarkan ID
func (s *AssessmentService) GetAssessmentByID(id uint) (*models.Assessment, error) {
	return s.AssessmentRepo.FindByID(id)
}

// GetAssessmentBySubmission mendapatkan penilaian berdasarkan ID pengajuan
func (s *AssessmentService) GetAssessmentBySubmission(submissionID uint) (*models.Assessment, error) {
	return s.AssessmentRepo.FindBySubmission(submissionID)
}

// UpdateAssessment memperbarui penilaian
func (s *AssessmentService) UpdateAssessment(assessment *models.Assessment) error {
	// Verifikasi penilaian yang ada
	existingAssessment, err := s.AssessmentRepo.FindByID(assessment.ID)
	if err != nil {
		return err
	}

	// Verifikasi adanya kiriman tersebut
	submission, err := s.SubmissionRepo.FindByID(existingAssessment.SubmissionID)
	if err != nil {
		return errors.New("submission not found")
	}

	// Verifikasi bahwa penugasan itu ada
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Validasi skor jika skor maksimal ditetapkan
	if assignment.MaxScore != nil && assessment.Score != nil {
		if *assessment.Score < 0 || *assessment.Score > *assignment.MaxScore {
			return errors.New("score must be between 0 and the maximum score")
		}
	}

	// Perbarui hanya bidang yang diizinkan
	existingAssessment.Score = assessment.Score
	existingAssessment.Feedback = assessment.Feedback
	existingAssessment.AssessedAt = time.Now()

	return s.AssessmentRepo.Update(existingAssessment)
}

// HapusPenilaian menghapus penilaian
func (s *AssessmentService) DeleteAssessment(id uint) error {
	return s.AssessmentRepo.Delete(id)
}
