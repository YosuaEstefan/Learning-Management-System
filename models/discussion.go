package models

import (
	"time"

	"gorm.io/gorm"

)

type Discussion struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null" json:"course_id"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
	Comments  []Comment `gorm:"foreignKey:DiscussionID" json:"comments,omitempty"`
}
