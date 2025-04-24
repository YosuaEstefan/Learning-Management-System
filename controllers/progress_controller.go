package controllers

import (
	"LMS/models"
	"LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

)

// ProgressController menangani permintaan kemajuan pembelajaran
type ProgressController struct {
	ProgressService *services.LearningProgressService
}

// NewProgressController membuat pengontrol progres baru
func NewProgressController(progressService *services.LearningProgressService) *ProgressController {
	return &ProgressController{
		ProgressService: progressService,
	}
}

// GradeRequest merupakan permintaan untuk menilai aktivitas siswa
type GradeRequest struct {
	StudentID    uint                `json:"student_id" binding:"required"`
	CourseID     uint                `json:"course_id" binding:"required"`
	ActivityType models.ProgressType `json:"activity_type" binding:"required"`
	ActivityID   uint                `json:"activity_id" binding:"required"`
	Score        float64             `json:"score" binding:"required"`
	MaxScore     float64             `json:"max_score" binding:"required"`
	Feedback     string              `json:"feedback"`
	Completed    bool                `json:"completed"`
}

// UpdateGradeRequest merepresentasikan permintaan untuk memperbarui nilai
type UpdateGradeRequest struct {
	Score     float64 `json:"score" binding:"required"`
	MaxScore  float64 `json:"max_score" binding:"required"`
	Feedback  string  `json:"feedback"`
	Completed bool    `json:"completed"`
}

// GradeActivity menangani penilaian aktivitas siswa
func (c *ProgressController) GradeActivity(ctx *gin.Context) {
	var request GradeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan ID grader dari konteks
	graderID, _ := ctx.Get("userID")

	err := c.ProgressService.CreateOrUpdateGrade(
		graderID.(uint),
		request.StudentID,
		request.CourseID,
		request.ActivityType,
		request.ActivityID,
		request.Score,
		request.MaxScore,
		request.Feedback,
		request.Completed,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Activity graded successfully",
	})
}

// GetProgressByID menangani mendapatkan kemajuan pembelajaran berdasarkan ID
func (c *ProgressController) GetProgressByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid progress ID"})
		return
	}

	progress, err := c.ProgressService.GetProgressByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Learning progress not found"})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// GetStudentProgress menangani mendapatkan semua kemajuan siswa dalam sebuah kursus
func (c *ProgressController) GetStudentProgress(ctx *gin.Context) {
	studentID, err := strconv.ParseUint(ctx.Param("student_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	progress, err := c.ProgressService.GetStudentProgress(uint(studentID), uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// GetCourseStudentsProgress menangani mendapatkan kemajuan untuk semua siswa dalam sebuah kursus
func (c *ProgressController) GetCourseStudentsProgress(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Dapatkan ID pemohon dari konteks
	requestorID, _ := ctx.Get("userID")

	progress, err := c.ProgressService.GetCourseStudentsProgress(uint(courseID), requestorID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// UpdateGradeByMentor menangani pembaruan nilai oleh mentor
func (c *ProgressController) UpdateGradeByMentor(ctx *gin.Context) {
	progressID, err := strconv.ParseUint(ctx.Param("progress_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid progress ID"})
		return
	}

	var request UpdateGradeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan ID mentor dari konteks
	mentorID, _ := ctx.Get("userID")

	err = c.ProgressService.UpdateGradeByMentor(
		uint(progressID),
		mentorID.(uint),
		request.Score,
		request.MaxScore,
		request.Feedback,
		request.Completed,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Grade updated successfully",
	})
}
