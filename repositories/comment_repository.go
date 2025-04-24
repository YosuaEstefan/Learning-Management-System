package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// CommentRepository handles database operations for comments
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// FindByID finds a comment by ID
func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	result := r.DB.First(&comment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, result.Error
	}
	return &comment, nil
}

// FindByDiscussion finds comments by discussion ID
func (r *CommentRepository) FindByDiscussion(discussionID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := r.DB.Where("discussion_id = ?", discussionID).Preload("User").Find(&comments)
	return comments, result.Error
}

// Create creates a new comment
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

// Update updates an existing comment
func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

// Delete deletes a comment
func (r *CommentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Comment{}, id).Error
}
