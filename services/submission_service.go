package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// SubmissionService menangani logika bisnis pengiriman tugas
type SubmissionService struct {
	SubmissionRepo *repositories.SubmissionRepository
	AssignmentRepo *repositories.AssignmentRepository
	EnrollmentRepo *repositories.EnrollmentRepository
	UserRepo       *repositories.UserRepository
}

// NewSubmissionService membuat layanan pengiriman baru
func NewSubmissionService(
	submissionRepo *repositories.SubmissionRepository,
	assignmentRepo *repositories.AssignmentRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
	userRepo *repositories.UserRepository,
) *SubmissionService {
	return &SubmissionService{
		SubmissionRepo: submissionRepo,
		AssignmentRepo: assignmentRepo,
		EnrollmentRepo: enrollmentRepo,
		UserRepo:       userRepo,
	}
}

// CreateSubmission membuat pengiriman baru untuk tugas
func (s *SubmissionService) CreateSubmission(submission *models.Submission, filePath string) error {
	// Verifikasi adanya penugasan
	assignment, err := s.AssignmentRepo.FindByID(submission.AssignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}
	// Verifikasi keberadaan siswa
	student, err := s.UserRepo.FindByID(submission.StudentID)
	if err != nil {
		return errors.New("student not found")
	}

	// Verifikasi bahwa siswa tersebut adalah siswa
	if student.Role != models.RoleStudent {
		return errors.New("only students can submit assignments")
	}

	// Verifikasi bahwa siswa telah terdaftar dalam kursus
	_, err = s.EnrollmentRepo.FindByUserAndCourse(submission.StudentID, assignment.CourseID)
	if err != nil {
		return errors.New("student is not enrolled in this course")
	}

	// Periksa tanggal jatuh tempo jika ditetapkan
	if assignment.DueDate != nil && time.Now().After(*assignment.DueDate) {
		return errors.New("assignment due date has passed")
	}

	// Periksa apakah siswa telah mengirimkan
	existingSubmission, err := s.SubmissionRepo.FindByAssignmentAndStudent(submission.AssignmentID, submission.StudentID)
	if err == nil && existingSubmission != nil {
		// Jika ada kiriman, perbarui
		existingSubmission.FilePath = filePath
		existingSubmission.SubmittedAt = time.Now()
		return s.SubmissionRepo.Update(existingSubmission)
	}
	// Mengatur jalur file dan tanggal pengiriman
	submission.FilePath = filePath
	submission.SubmittedAt = time.Now()

	// Buat kiriman baru
	return s.SubmissionRepo.Create(submission)
}

// GetSubmissionByID berdasrkan id
func (s *SubmissionService) GetSubmissionByID(id uint) (*models.Submission, error) {
	return s.SubmissionRepo.FindByID(id)
}

// GetSubmissionsByAssignment berdasrkan assignment ID
func (s *SubmissionService) GetSubmissionsByAssignment(assignmentID uint) ([]models.Submission, error) {
	return s.SubmissionRepo.FindByAssignment(assignmentID)
}

// GetSubmissionsByStudent berdasrkan student ID
func (s *SubmissionService) GetSubmissionsByStudent(studentID uint) ([]models.Submission, error) {
	return s.SubmissionRepo.FindByStudent(studentID)
}

// DeleteSubmission untuk submission
func (s *SubmissionService) DeleteSubmission(id uint) error {
	return s.SubmissionRepo.Delete(id)
}
