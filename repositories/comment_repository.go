package repositories

import (
	"LMS/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetAllByDiscussionID(discussionID uint) ([]models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetAllByDiscussionID(discussionID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("discussion_id = ?", discussionID).Find(&comments).Error
	return comments, err
}
