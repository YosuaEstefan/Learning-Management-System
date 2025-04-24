package models

import (
	"time"

	"gorm.io/gorm"

)

type ProgressType string

const (
	ProgressTypeMaterial   ProgressType = "material"
	ProgressTypeAssignment ProgressType = "assignment"
	ProgressTypeQuiz       ProgressType = "quiz"
	ProgressTypeDiscussion ProgressType = "discussion"
)

type LearningProgress struct {
	gorm.Model
	ID            uint         `gorm:"primaryKey" json:"id"`
	UserID        uint         `gorm:"not null" json:"user_id"`
	User          User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID      uint         `gorm:"not null" json:"course_id"`
	Course        Course       `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	ActivityType  ProgressType `gorm:"type:enum('material','assignment','quiz','discussion');not null" json:"activity_type"`
	ActivityID    uint         `gorm:"not null" json:"activity_id"`
	ActivityTitle string       `gorm:"-" json:"activity_title,omitempty"`
	Score         *float64     `json:"score"`
	MaxScore      *float64     `json:"max_score"`
	Feedback      string       `gorm:"type:text" json:"feedback"`
	GradedBy      uint         `json:"graded_by"`
	Grader        User         `gorm:"foreignKey:GradedBy" json:"grader,omitempty"`
	GraderName    string       `gorm:"-" json:"grader_name,omitempty"`
	Completed     bool         `gorm:"default:false" json:"completed"`
	CompletedAt   *time.Time   `json:"completed_at"`
	CreatedAt     time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
