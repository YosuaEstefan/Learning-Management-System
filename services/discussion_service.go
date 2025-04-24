package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// DiscussionService menangani logika bisnis diskusi
type DiscussionService struct {
	DiscussionRepo *repositories.DiscussionRepository
	CourseRepo     *repositories.CourseRepository
	UserRepo       *repositories.UserRepository
	EnrollmentRepo *repositories.EnrollmentRepository
}

// NewDiscussionService membuat layanan diskusi baru
func NewDiscussionService(
	discussionRepo *repositories.DiscussionRepository,
	courseRepo *repositories.CourseRepository,
	userRepo *repositories.UserRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
) *DiscussionService {
	return &DiscussionService{
		DiscussionRepo: discussionRepo,
		CourseRepo:     courseRepo,
		UserRepo:       userRepo,
		EnrollmentRepo: enrollmentRepo,
	}
}

// CreateDiscussion membuat diskusi baru
func (s *DiscussionService) CreateDiscussion(discussion *models.Discussion) error {
	// Verifikasi keberadaan kursus
	_, err := s.CourseRepo.FindByID(discussion.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Verifikasi keberadaan pengguna
	user, err := s.UserRepo.FindByID(discussion.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Jika pengguna adalah siswa, verifikasi pendaftaran
	if user.Role == models.RoleStudent {
		_, err = s.EnrollmentRepo.FindByUserAndCourse(discussion.UserID, discussion.CourseID)
		if err != nil {
			return errors.New("student is not enrolled in this course")
		}
	}

	// Tetapkan tanggal pembuatan
	discussion.CreatedAt = time.Now()

	// Buat diskusi
	return s.DiscussionRepo.Create(discussion)
}

// GetDiscussionByID mendapatkan diskusi berdasarkan ID
func (s *DiscussionService) GetDiscussionByID(id uint) (*models.Discussion, error) {
	return s.DiscussionRepo.FindByID(id)
}

// GetDiscussionsByCourse mendapatkan diskusi berdasarkan ID kursus
func (s *DiscussionService) GetDiscussionsByCourse(courseID uint) ([]models.Discussion, error) {
	return s.DiscussionRepo.FindByCourse(courseID)
}

// UpdateDiscussion memperbarui diskusi
func (s *DiscussionService) UpdateDiscussion(discussion *models.Discussion) error {
	// Verifikasi adanya diskusi
	existingDiscussion, err := s.DiscussionRepo.FindByID(discussion.ID)
	if err != nil {
		return err
	}

	// Verifikasi pengguna sebagai pemilik diskusi
	if existingDiscussion.UserID != discussion.UserID {
		return errors.New("user does not own this discussion")
	}

	// Perbarui hanya bidang yang diizinkan
	existingDiscussion.Title = discussion.Title
	existingDiscussion.Content = discussion.Content

	return s.DiscussionRepo.Update(existingDiscussion)
}

// DeleteDiscussion menghapus diskusi
func (s *DiscussionService) DeleteDiscussion(id uint, userID uint, isAdmin bool) error {
	// Verifikasi adanya diskusi
	discussion, err := s.DiscussionRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verifikasi pengguna yang memiliki diskusi atau admin
	if !isAdmin && discussion.UserID != userID {
		return errors.New("user does not own this discussion")
	}

	return s.DiscussionRepo.Delete(id)
}
