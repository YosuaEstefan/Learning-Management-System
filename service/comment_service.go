package services

import (
	"LMS/models"
	"LMS/repositories"
)

var CommentServiceInstance CommentService

type CommentService interface {
	CreateComment(comment *models.Comment) error
	GetCommentsByDiscussionID(discussionID uint) ([]models.Comment, error)
}

type commentService struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &commentService{repo}
}

func (s *commentService) CreateComment(comment *models.Comment) error {
	return s.repo.Create(comment)
}

func (s *commentService) GetCommentsByDiscussionID(discussionID uint) ([]models.Comment, error) {
	return s.repo.GetAllByDiscussionID(discussionID)
}
