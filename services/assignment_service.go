package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// AssignmentService menangani logika bisnis penugasan
type AssignmentService struct {
	AssignmentRepo *repositories.AssignmentRepository
	CourseRepo     *repositories.CourseRepository
}

// NewAssignmentService membuat layanan penugasan baru
func NewAssignmentService(assignmentRepo *repositories.AssignmentRepository, courseRepo *repositories.CourseRepository) *AssignmentService {
	return &AssignmentService{
		AssignmentRepo: assignmentRepo,
		CourseRepo:     courseRepo,
	}
}

// CreateAssignment membuat penugasan baru
func (s *AssignmentService) CreateAssignment(assignment *models.Assignment) error {
	// Verifikasi keberadaan kursus
	_, err := s.CourseRepo.FindByID(assignment.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Tetapkan tanggal pembuatan
	assignment.CreatedAt = time.Now()

	// Buat tugas
	return s.AssignmentRepo.Create(assignment)
}

// GetAssignmentByID mendapatkan penugasan dengan ID
func (s *AssignmentService) GetAssignmentByID(id uint) (*models.Assignment, error) {
	return s.AssignmentRepo.FindByID(id)
}

// GetAssignmentsByCourse mendapatkan tugas berdasarkan ID kursus
func (s *AssignmentService) GetAssignmentsByCourse(courseID uint) ([]models.Assignment, error) {
	return s.AssignmentRepo.FindByCourse(courseID)
}

// UpdateAssignment memperbarui penugasan
func (s *AssignmentService) UpdateAssignment(assignment *models.Assignment) error {
	// Verifikasi bahwa penugasan itu ada
	existingAssignment, err := s.AssignmentRepo.FindByID(assignment.ID)
	if err != nil {
		return err
	}

	// Perbarui hanya bidang yang diizinkan
	existingAssignment.Title = assignment.Title
	existingAssignment.Description = assignment.Description
	existingAssignment.DueDate = assignment.DueDate
	existingAssignment.MaxScore = assignment.MaxScore

	return s.AssignmentRepo.Update(existingAssignment)
}

// DeleteAssignment menghapus penugasan
func (s *AssignmentService) DeleteAssignment(id uint) error {
	return s.AssignmentRepo.Delete(id)
}
