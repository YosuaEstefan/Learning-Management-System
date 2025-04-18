package repositories

import (
	"LMS/models"

	"gorm.io/gorm"
)

type DiscussionRepository interface {
	Create(d *models.Discussion) error
	GetByID(id uint) (*models.Discussion, error)
	GetAllByCourseID(courseID uint) ([]models.Discussion, error)
}

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &discussionRepository{db}
}

func (r *discussionRepository) Create(d *models.Discussion) error {
	return r.db.Create(d).Error
}

func (r *discussionRepository) GetByID(id uint) (*models.Discussion, error) {
	var d models.Discussion
	err := r.db.Preload("Comments").First(&d, id).Error
	return &d, err
}

func (r *discussionRepository) GetAllByCourseID(courseID uint) ([]models.Discussion, error) {
	var discussions []models.Discussion
	err := r.db.Where("course_id = ?", courseID).Find(&discussions).Error
	return discussions, err
}
