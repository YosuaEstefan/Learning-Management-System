// services/discussion_service.go
package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// DiscussionService handles discussion business logic
type DiscussionService struct {
	DiscussionRepo *repositories.DiscussionRepository
	CourseRepo     *repositories.CourseRepository
	UserRepo       *repositories.UserRepository
	EnrollmentRepo *repositories.EnrollmentRepository
}

// NewDiscussionService creates a new discussion service
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

// CreateDiscussion creates a new discussion
func (s *DiscussionService) CreateDiscussion(discussion *models.Discussion) error {
	// Verify the course exists
	_, err := s.CourseRepo.FindByID(discussion.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Verify the user exists
	user, err := s.UserRepo.FindByID(discussion.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// If user is a student, verify enrollment
	if user.Role == models.RoleStudent {
		_, err = s.EnrollmentRepo.FindByUserAndCourse(discussion.UserID, discussion.CourseID)
		if err != nil {
			return errors.New("student is not enrolled in this course")
		}
	}

	// Set creation date
	discussion.CreatedAt = time.Now()

	// Create the discussion
	return s.DiscussionRepo.Create(discussion)
}

// GetDiscussionByID gets a discussion by ID
func (s *DiscussionService) GetDiscussionByID(id uint) (*models.Discussion, error) {
	return s.DiscussionRepo.FindByID(id)
}

// GetDiscussionsByCourse gets discussions by course ID
func (s *DiscussionService) GetDiscussionsByCourse(courseID uint) ([]models.Discussion, error) {
	return s.DiscussionRepo.FindByCourse(courseID)
}

// UpdateDiscussion updates a discussion
func (s *DiscussionService) UpdateDiscussion(discussion *models.Discussion) error {
	// Verify the discussion exists
	existingDiscussion, err := s.DiscussionRepo.FindByID(discussion.ID)
	if err != nil {
		return err
	}

	// Verify the user owns the discussion
	if existingDiscussion.UserID != discussion.UserID {
		return errors.New("user does not own this discussion")
	}

	// Update only allowed fields
	existingDiscussion.Title = discussion.Title
	existingDiscussion.Content = discussion.Content

	return s.DiscussionRepo.Update(existingDiscussion)
}

// DeleteDiscussion deletes a discussion
func (s *DiscussionService) DeleteDiscussion(id uint, userID uint, isAdmin bool) error {
	// Verify the discussion exists
	discussion, err := s.DiscussionRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify the user owns the discussion or is an admin
	if !isAdmin && discussion.UserID != userID {
		return errors.New("user does not own this discussion")
	}

	return s.DiscussionRepo.Delete(id)
}
