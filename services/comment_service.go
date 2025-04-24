package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// CommentService menangani logika bisnis komentar
type CommentService struct {
	CommentRepo    *repositories.CommentRepository
	DiscussionRepo *repositories.DiscussionRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
	EnrollmentRepo *repositories.EnrollmentRepository
}

// NewCommentService membuat layanan komentar baru
func NewCommentService(
	commentRepo *repositories.CommentRepository,
	discussionRepo *repositories.DiscussionRepository,
	userRepo *repositories.UserRepository,
	courseRepo *repositories.CourseRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
) *CommentService {
	return &CommentService{
		CommentRepo:    commentRepo,
		DiscussionRepo: discussionRepo,
		UserRepo:       userRepo,
		CourseRepo:     courseRepo,
		EnrollmentRepo: enrollmentRepo,
	}
}

// CreateComment membuat komentar baru
func (s *CommentService) CreateComment(comment *models.Comment) error {
	// Verifikasi adanya diskusi
	discussion, err := s.DiscussionRepo.FindByID(comment.DiscussionID)
	if err != nil {
		return errors.New("discussion not found")
	}

	// Verifikasi keberadaan pengguna
	user, err := s.UserRepo.FindByID(comment.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Jika pengguna adalah siswa, verifikasi pendaftaran
	if user.Role == models.RoleStudent {
		_, err = s.EnrollmentRepo.FindByUserAndCourse(comment.UserID, discussion.CourseID)
		if err != nil {
			return errors.New("student is not enrolled in this course")
		}
	}

	// Tetapkan tanggal pembuatan
	comment.CreatedAt = time.Now()
	// Buat komentar
	return s.CommentRepo.Create(comment)
}

// GetCommentByID mendapatkan komentar berdasarkan ID
func (s *CommentService) GetCommentByID(id uint) (*models.Comment, error) {
	return s.CommentRepo.FindByID(id)
}

// GetCommentsByDiscussion mendapatkan komentar berdasarkan ID diskusi
func (s *CommentService) GetCommentsByDiscussion(discussionID uint) ([]models.Comment, error) {
	return s.CommentRepo.FindByDiscussion(discussionID)
}

// UpdateComment memperbarui komentar
func (s *CommentService) UpdateComment(comment *models.Comment) error {
	// Verifikasi bahwa komentar tersebut ada
	existingComment, err := s.CommentRepo.FindByID(comment.ID)
	if err != nil {
		return err
	}

	// Verifikasi pengguna yang memiliki komentar
	if existingComment.UserID != comment.UserID {
		return errors.New("user does not own this comment")
	}
	// Perbarui hanya bidang yang diizinkan
	existingComment.Content = comment.Content

	return s.CommentRepo.Update(existingComment)
}

// DeleteComment menghapus komentar
func (s *CommentService) DeleteComment(id uint, userID uint, isAdmin bool) error {
	// Verifikasi bahwa komentar tersebut ada
	comment, err := s.CommentRepo.FindByID(id)
	if err != nil {
		return err
	}
	// Verifikasi pengguna yang memiliki komentar atau admin
	if !isAdmin && comment.UserID != userID {
		return errors.New("user does not own this comment")
	}

	return s.CommentRepo.Delete(id)
}
