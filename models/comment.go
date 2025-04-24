package models

import (
	"time"

	"gorm.io/gorm"

)

type Comment struct {
	gorm.Model
	ID           uint       `gorm:"primaryKey" json:"id"`
	DiscussionID uint       `gorm:"not null" json:"discussion_id"`
	Discussion   Discussion `gorm:"foreignKey:DiscussionID" json:"discussion,omitempty"`
	UserID       uint       `gorm:"not null" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content      string     `gorm:"type:text" json:"content"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}
