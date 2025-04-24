package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// CommentRepository menangani operasi basis data untuk komentar
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository membuat repositori komentar baru
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// FindByID menemukan komentar berdasarkan ID
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

// FindByDiscussion menemukan komentar berdasarkan ID diskusi
func (r *CommentRepository) FindByDiscussion(discussionID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := r.DB.Where("discussion_id = ?", discussionID).Preload("User").Find(&comments)
	return comments, result.Error
}

// Buat membuat komentar baru
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

// Perbarui memperbarui komentar yang ada
func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

// Hapus menghapus komentar
func (r *CommentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Comment{}, id).Error
}
