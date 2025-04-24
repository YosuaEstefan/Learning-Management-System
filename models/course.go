package models

import (
	"time"

	"gorm.io/gorm"

)

type Course struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey" json:"id"`
	Title       string       `gorm:"size:255;not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	MentorID    uint         `gorm:"not null" json:"mentor_id"`
	Mentor      User         `gorm:"foreignKey:MentorID" json:"mentor,omitempty"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
	Materials   []Material   `gorm:"foreignKey:CourseID" json:"materials,omitempty"`
	Assignments []Assignment `gorm:"foreignKey:CourseID" json:"assignments,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID" json:"enrollments,omitempty"`
	Discussions []Discussion `gorm:"foreignKey:CourseID" json:"discussions,omitempty"`
}
