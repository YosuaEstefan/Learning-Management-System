package models

import (
	"time"

	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	CourseID       uint           `gorm:"not null" json:"course_id"`
	EnrollmentDate time.Time      `gorm:"autoCreateTime" json:"enrollment_date"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
