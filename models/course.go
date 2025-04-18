package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	MentorID    uint           `gorm:"not null" json:"mentor_id"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// Relasi
	Mentor User `gorm:"foreignKey:MentorID" json:"mentor"`
}
