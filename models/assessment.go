package models

import (
	"time"

	"gorm.io/gorm"

)

type Assessment struct {
	gorm.Model
	ID           uint       `gorm:"primaryKey" json:"id"`
	SubmissionID uint       `gorm:"not null;uniqueIndex" json:"submission_id"`
	Submission   Submission `gorm:"foreignKey:SubmissionID" json:"submission,omitempty"`
	Score        *int       `json:"score"`
	Feedback     string     `gorm:"type:text" json:"feedback"`
	AssessedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"assessed_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}
