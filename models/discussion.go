package models

import (
	"time"

	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null" json:"course_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Comments  []Comment `json:"comments,omitempty"`
}
