package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// DiscussionRepository handles database operations for discussions
type DiscussionRepository struct {
	DB *gorm.DB
}

// NewDiscussionRepository creates a new discussion repository
func NewDiscussionRepository(db *gorm.DB) *DiscussionRepository {
	return &DiscussionRepository{DB: db}
}

// FindByID finds a discussion by ID
func (r *DiscussionRepository) FindByID(id uint) (*models.Discussion, error) {
	var discussion models.Discussion
	result := r.DB.First(&discussion, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("discussion not found")
		}
		return nil, result.Error
	}
	return &discussion, nil
}

// FindByCourse finds discussions by course ID
func (r *DiscussionRepository) FindByCourse(courseID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&discussions)
	return discussions, result.Error
}

// FindByUser finds discussions by user ID
func (r *DiscussionRepository) FindByUser(userID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	result := r.DB.Where("user_id = ?", userID).Preload("Course").Find(&discussions)
	return discussions, result.Error
}

// Create creates a new discussion
func (r *DiscussionRepository) Create(discussion *models.Discussion) error {
	return r.DB.Create(discussion).Error
}

// Update updates an existing discussion
func (r *DiscussionRepository) Update(discussion *models.Discussion) error {
	return r.DB.Save(discussion).Error
}

// Delete deletes a discussion
func (r *DiscussionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Discussion{}, id).Error
}
