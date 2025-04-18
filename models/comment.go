package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey" json:"id"`
	DiscussionID uint      `gorm:"not null" json:"discussion_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	Content      string    `gorm:"type:text" json:"content"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
