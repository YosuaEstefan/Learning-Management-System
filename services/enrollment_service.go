package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// EnrollmentService menangani logika bisnis pendaftaran
type EnrollmentService struct {
	EnrollmentRepo *repositories.EnrollmentRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
}

// NewEnrollmentService membuat layanan pendaftaran baru
func NewEnrollmentService(enrollmentRepo *repositories.EnrollmentRepository, userRepo *repositories.UserRepository, courseRepo *repositories.CourseRepository) *EnrollmentService {
	return &EnrollmentService{
		EnrollmentRepo: enrollmentRepo,
		UserRepo:       userRepo,
		CourseRepo:     courseRepo,
	}
}

// CreateEnrollment membuat pendaftaran baru
func (s *EnrollmentService) CreateEnrollment(enrollment *models.Enrollment) error {
	// Verifikasi keberadaan pengguna
	user, err := s.UserRepo.FindByID(enrollment.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verifikasi pengguna adalah siswa
	if user.Role != models.RoleStudent {
		return errors.New("only students can enroll in courses")
	}

	// Verifikasi keberadaan kursus
	_, err = s.CourseRepo.FindByID(enrollment.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Periksa apakah siswa sudah terdaftar
	_, err = s.EnrollmentRepo.FindByUserAndCourse(enrollment.UserID, enrollment.CourseID)
	if err == nil {
		return errors.New("student is already enrolled in this course")
	}

	// Tetapkan tanggal pendaftaran
	enrollment.EnrollmentDate = time.Now()

	// Buat pendaftaran
	return s.EnrollmentRepo.Create(enrollment)
}

// GetEnrollmentByID mendapatkan pendaftaran dengan ID
func (s *EnrollmentService) GetEnrollmentByID(id uint) (*models.Enrollment, error) {
	return s.EnrollmentRepo.FindByID(id)
}

// GetEnrollmentsByUser mendapatkan pendaftaran berdasarkan ID pengguna
func (s *EnrollmentService) GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error) {
	return s.EnrollmentRepo.FindByUser(userID)
}

// GetEnrollmentsByCourse mendapatkan pendaftaran berdasarkan ID kursus
func (s *EnrollmentService) GetEnrollmentsByCourse(courseID uint) ([]models.Enrollment, error) {
	return s.EnrollmentRepo.FindByCourse(courseID)
}

// HapusPendaftaran menghapus pendaftaran
func (s *EnrollmentService) DeleteEnrollment(id uint) error {
	return s.EnrollmentRepo.Delete(id)
}

// IsStudentEnrolled memeriksa apakah seorang siswa terdaftar dalam sebuah kursus
func (s *EnrollmentService) IsStudentEnrolled(userID, courseID uint) (bool, error) {
	_, err := s.EnrollmentRepo.FindByUserAndCourse(userID, courseID)
	if err != nil {
		return false, nil
	}
	return true, nil
}
