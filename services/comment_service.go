package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// CommentService handles comment business logic
type CommentService struct {
	CommentRepo    *repositories.CommentRepository
	DiscussionRepo *repositories.DiscussionRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
	EnrollmentRepo *repositories.EnrollmentRepository
}

// NewCommentService creates a new comment service
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

// CreateComment creates a new comment
func (s *CommentService) CreateComment(comment *models.Comment) error {
	// Verify the discussion exists
	discussion, err := s.DiscussionRepo.FindByID(comment.DiscussionID)
	if err != nil {
		return errors.New("discussion not found")
	}

	// Verify the user exists
	user, err := s.UserRepo.FindByID(comment.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// If user is a student, verify enrollment
	if user.Role == models.RoleStudent {
		_, err = s.EnrollmentRepo.FindByUserAndCourse(comment.UserID, discussion.CourseID)
		if err != nil {
			return errors.New("student is not enrolled in this course")
		}
	}

	// Set creation date
	comment.CreatedAt = time.Now()

	// Create the comment
	return s.CommentRepo.Create(comment)
}

// GetCommentByID gets a comment by ID
func (s *CommentService) GetCommentByID(id uint) (*models.Comment, error) {
	return s.CommentRepo.FindByID(id)
}

// GetCommentsByDiscussion gets comments by discussion ID
func (s *CommentService) GetCommentsByDiscussion(discussionID uint) ([]models.Comment, error) {
	return s.CommentRepo.FindByDiscussion(discussionID)
}

// UpdateComment updates a comment
func (s *CommentService) UpdateComment(comment *models.Comment) error {
	// Verify the comment exists
	existingComment, err := s.CommentRepo.FindByID(comment.ID)
	if err != nil {
		return err
	}

	// Verify the user owns the comment
	if existingComment.UserID != comment.UserID {
		return errors.New("user does not own this comment")
	}

	// Update only allowed fields
	existingComment.Content = comment.Content

	return s.CommentRepo.Update(existingComment)
}

// DeleteComment deletes a comment
func (s *CommentService) DeleteComment(id uint, userID uint, isAdmin bool) error {
	// Verify the comment exists
	comment, err := s.CommentRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify the user owns the comment or is an admin
	if !isAdmin && comment.UserID != userID {
		return errors.New("user does not own this comment")
	}

	return s.CommentRepo.Delete(id)
}
