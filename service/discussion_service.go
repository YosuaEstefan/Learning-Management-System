package services

import (
	"LMS/models"
	"LMS/repositories"
)

var DiscussionServiceInstance DiscussionService

type DiscussionService interface {
	CreateDiscussion(d *models.Discussion) error
	GetDiscussionByID(id uint) (*models.Discussion, error)
	GetDiscussionsByCourseID(courseID uint) ([]models.Discussion, error)
}

type discussionService struct {
	repo repositories.DiscussionRepository
}

func NewDiscussionService(repo repositories.DiscussionRepository) DiscussionService {
	return &discussionService{repo}
}

func (s *discussionService) CreateDiscussion(d *models.Discussion) error {
	return s.repo.Create(d)
}

func (s *discussionService) GetDiscussionByID(id uint) (*models.Discussion, error) {
	return s.repo.GetByID(id)
}

func (s *discussionService) GetDiscussionsByCourseID(courseID uint) ([]models.Discussion, error) {
	return s.repo.GetAllByCourseID(courseID)
}
