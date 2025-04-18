package models

import (
	"time"

	"gorm.io/gorm"
)

type Assessment struct {
	gorm.Model
	ID           uint           `gorm:"primaryKey" json:"id"`
	SubmissionID uint           `gorm:"not null" json:"submission_id"`
	Score        int            `json:"score"`
	Feedback     string         `gorm:"type:text" json:"feedback"`
	AssessedAt   time.Time      `gorm:"autoCreateTime" json:"assessed_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
