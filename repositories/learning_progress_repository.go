// repositories/learning_progress_repository.go
package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"
)

// LearningProgressRepository handles database operations for learning progress
type LearningProgressRepository struct {
	DB *gorm.DB
}

// NewLearningProgressRepository creates a new learning progress repository
func NewLearningProgressRepository(db *gorm.DB) *LearningProgressRepository {
	return &LearningProgressRepository{DB: db}
}

// FindByID finds a learning progress by ID
func (r *LearningProgressRepository) FindByID(id uint) (*models.LearningProgress, error) {
	var progress models.LearningProgress
	result := r.DB.Preload("User").Preload("Course").Preload("Grader").First(&progress, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("learning progress not found")
		}
		return nil, result.Error
	}
	return &progress, nil
}

// GetStudentProgress gets all progress for a student in a course
func (r *LearningProgressRepository) GetStudentProgress(studentID, courseID uint) ([]models.LearningProgress, error) {
	var progressItems []models.LearningProgress

	result := r.DB.Where("user_id = ? AND course_id = ?", studentID, courseID)
	result = result.Preload("User").Preload("Course").Preload("Grader")
	result = result.Order("activity_type, activity_id")
	result = result.Find(&progressItems)

	if result.Error != nil {
		return nil, result.Error
	}

	// Populate computed fields that are not directly stored in DB
	for i := range progressItems {
		if progressItems[i].Grader.ID > 0 {
			progressItems[i].GraderName = progressItems[i].Grader.Name
		}

		// Load activity title based on activity type
		switch progressItems[i].ActivityType {
		case models.ProgressTypeAssignment:
			var assignment models.Assignment
			if err := r.DB.First(&assignment, progressItems[i].ActivityID).Error; err == nil {
				progressItems[i].ActivityTitle = assignment.Title
			} else {
				progressItems[i].ActivityTitle = "Unknown Assignment"
			}
		case models.ProgressTypeMaterial:
			var material models.Material
			if err := r.DB.First(&material, progressItems[i].ActivityID).Error; err == nil {
				progressItems[i].ActivityTitle = material.Title
			} else {
				progressItems[i].ActivityTitle = "Unknown Material"
			}
		case models.ProgressTypeQuiz:
			progressItems[i].ActivityTitle = "Quiz " + string(progressItems[i].ActivityID)
		case models.ProgressTypeDiscussion:
			var discussion models.Discussion
			if err := r.DB.First(&discussion, progressItems[i].ActivityID).Error; err == nil {
				progressItems[i].ActivityTitle = discussion.Title
			} else {
				progressItems[i].ActivityTitle = "Unknown Discussion"
			}
		}
	}

	return progressItems, nil
}

// GetCourseStudentsProgress gets progress for all students in a course
func (r *LearningProgressRepository) GetCourseStudentsProgress(courseID uint) (map[uint][]models.LearningProgress, error) {
	// Get all students enrolled in the course
	var enrollments []models.Enrollment
	if err := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&enrollments).Error; err != nil {
		return nil, err
	}

	result := make(map[uint][]models.LearningProgress)

	// Get progress for each student
	for _, enrollment := range enrollments {
		if enrollment.User.Role != models.RoleStudent {
			continue
		}

		studentProgress, err := r.GetStudentProgress(enrollment.UserID, courseID)
		if err != nil {
			continue
		}

		result[enrollment.UserID] = studentProgress
	}

	return result, nil
}

// Create creates a new learning progress
func (r *LearningProgressRepository) Create(progress *models.LearningProgress) error {
	return r.DB.Create(progress).Error
}

// Update updates an existing learning progress
func (r *LearningProgressRepository) Update(progress *models.LearningProgress) error {
	return r.DB.Model(&models.LearningProgress{}).Where("id = ?", progress.ID).Updates(progress).Error
}

// Delete deletes a learning progress
func (r *LearningProgressRepository) Delete(id uint) error {
	return r.DB.Delete(&models.LearningProgress{}, id).Error
}

// FindByActivityAndStudent finds progress by activity type, activity ID, and student ID
func (r *LearningProgressRepository) FindByActivityAndStudent(
	activityType models.ProgressType,
	activityID uint,
	studentID uint,
	courseID uint,
) (*models.LearningProgress, error) {
	var progress models.LearningProgress

	result := r.DB.Where(
		"activity_type = ? AND activity_id = ? AND user_id = ? AND course_id = ?",
		activityType, activityID, studentID, courseID,
	).First(&progress)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("progress not found")
		}
		return nil, result.Error
	}

	return &progress, nil
}
