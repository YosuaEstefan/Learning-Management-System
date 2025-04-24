package models

import (
	"time"

	"gorm.io/gorm"

)

type Submission struct {
	gorm.Model
	ID           uint        `gorm:"primaryKey" json:"id"`
	AssignmentID uint        `gorm:"not null" json:"assignment_id"`
	Assignment   Assignment  `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	StudentID    uint        `gorm:"not null" json:"student_id"`
	Student      User        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	FilePath     string      `gorm:"size:255;not null" json:"file_path"`
	SubmittedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"submitted_at"`
	UpdatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
	Assessment   *Assessment `gorm:"foreignKey:SubmissionID" json:"assessment,omitempty"`
}
