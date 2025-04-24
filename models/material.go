package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey" json:"id"`
	CourseID   uint      `gorm:"not null" json:"course_id"`
	Course     Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	FilePath   string    `gorm:"size:255;not null" json:"file_path"`
	UploadedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"uploaded_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}

// type Material struct {
// 	gorm.Model
// 	ID         uint           `gorm:"primaryKey" json:"id"`
// 	CourseID   uint           `gorm:"not null" json:"course_id"`
// 	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
// 	FilePath   string         `gorm:"type:varchar(255);not null" json:"file_path"`
// 	UploadedAt time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
// 	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
// }
