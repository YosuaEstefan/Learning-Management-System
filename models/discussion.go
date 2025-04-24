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

	// Relations
	Comments []Comment `gorm:"foreignKey:DiscussionID" json:"comments,omitempty"`
}

// type Discussion struct {
// 	gorm.Model
// 	ID        uint      `gorm:"primaryKey" json:"id"`
// 	CourseID  uint      `json:"course_id"`
// 	UserID    uint      `json:"user_id"`
// 	Title     string    `json:"title"`
// 	Content   string    `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// 	Comments  []Comment `json:"comments,omitempty"`
// }
