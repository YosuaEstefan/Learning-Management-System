package models

import (
	"time"

	"gorm.io/gorm"

)

type Enrollment struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID       uint      `gorm:"not null" json:"course_id"`
	Course         Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	EnrollmentDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"enrollment_date"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}
