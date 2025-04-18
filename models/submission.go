package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	ID           uint           `gorm:"primaryKey" json:"id"`
	AssignmentID uint           `gorm:"not null" json:"assignment_id"`
	StudentID    uint           `gorm:"not null" json:"student_id"`
	FilePath     string         `gorm:"type:text";not null" json:"file_path"`
	SubmittedAt  time.Time      `gorm:"autoCreateTime" json:"submitted_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
