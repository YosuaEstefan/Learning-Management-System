package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	gorm.Model
	ID         uint           `gorm:"primaryKey" json:"id"`
	CourseID   uint           `gorm:"not null" json:"course_id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	FilePath   string         `gorm:"type:text";not null" json:"file_path"`
	UploadedAt time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
