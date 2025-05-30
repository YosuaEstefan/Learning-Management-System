package models

import (
	"time"

	"gorm.io/gorm"

)

type Assignment struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey" json:"id"`
	CourseID    uint         `gorm:"not null" json:"course_id"`
	Course      Course       `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Title       string       `gorm:"size:255;not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	DueDate     *time.Time   `json:"due_date"`
	MaxScore    *int         `json:"max_score"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions,omitempty"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}
