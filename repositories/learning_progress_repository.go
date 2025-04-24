package repositories

import (
	"LMS/models"
	"errors"

	"gorm.io/gorm"

)

// LearningProgressRepository menangani operasi basis data untuk kemajuan pembelajaran
type LearningProgressRepository struct {
	DB *gorm.DB
}

// NewLearningProgressRepository membuat repositori kemajuan pembelajaran baru
func NewLearningProgressRepository(db *gorm.DB) *LearningProgressRepository {
	return &LearningProgressRepository{DB: db}
}

// FindByID menemukan kemajuan pembelajaran berdasarkan ID
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

// GetStudentProgress mendapatkan semua kemajuan siswa dalam sebuah kursus
func (r *LearningProgressRepository) GetStudentProgress(studentID, courseID uint) ([]models.LearningProgress, error) {
	var progressItems []models.LearningProgress

	result := r.DB.Where("user_id = ? AND course_id = ?", studentID, courseID)
	result = result.Preload("User").Preload("Course").Preload("Grader")
	result = result.Order("activity_type, activity_id")
	result = result.Find(&progressItems)

	if result.Error != nil {
		return nil, result.Error
	}

	// Mengisi bidang yang dihitung yang tidak disimpan secara langsung di DB
	for i := range progressItems {
		if progressItems[i].Grader.ID > 0 {
			progressItems[i].GraderName = progressItems[i].Grader.Name
		}

		// Memuat judul aktivitas berdasarkan jenis aktivitas
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

// GetCourseStudentsProgress mendapatkan kemajuan untuk semua siswa dalam sebuah kursus
func (r *LearningProgressRepository) GetCourseStudentsProgress(courseID uint) (map[uint][]models.LearningProgress, error) {
	// Dapatkan semua siswa yang terdaftar dalam kursus
	var enrollments []models.Enrollment
	if err := r.DB.Where("course_id = ?", courseID).Preload("User").Find(&enrollments).Error; err != nil {
		return nil, err
	}

	result := make(map[uint][]models.LearningProgress)

	// Dapatkan kemajuan untuk setiap siswa
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

// Membuat menciptakan kemajuan pembelajaran baru
func (r *LearningProgressRepository) Create(progress *models.LearningProgress) error {
	return r.DB.Create(progress).Error
}

// Memperbarui kemajuan pembelajaran yang sudah ada
func (r *LearningProgressRepository) Update(progress *models.LearningProgress) error {
	return r.DB.Model(&models.LearningProgress{}).Where("id = ?", progress.ID).Updates(progress).Error
}

// Hapus menghapus kemajuan pembelajaran
func (r *LearningProgressRepository) Delete(id uint) error {
	return r.DB.Delete(&models.LearningProgress{}, id).Error
}

// FindByActivityAndStudent menemukan kemajuan berdasarkan jenis aktivitas, ID aktivitas, dan ID siswa
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
