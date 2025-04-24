// controllers/progress_controller.go
package controllers

import (
	"LMS/models"
	"LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProgressController handles learning progress requests
type ProgressController struct {
	ProgressService *services.LearningProgressService
}

// NewProgressController creates a new progress controller
func NewProgressController(progressService *services.LearningProgressService) *ProgressController {
	return &ProgressController{
		ProgressService: progressService,
	}
}

// GradeRequest represents a request to grade a student's activity
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

// UpdateGradeRequest represents a request to update a grade
type UpdateGradeRequest struct {
	Score     float64 `json:"score" binding:"required"`
	MaxScore  float64 `json:"max_score" binding:"required"`
	Feedback  string  `json:"feedback"`
	Completed bool    `json:"completed"`
}

// GradeActivity handles the grading of a student's activity
func (c *ProgressController) GradeActivity(ctx *gin.Context) {
	var request GradeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get grader ID from context
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

// GetProgressByID handles getting a learning progress by ID
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

// GetStudentProgress handles getting all progress for a student in a course
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

// GetCourseStudentsProgress handles getting progress for all students in a course
func (c *ProgressController) GetCourseStudentsProgress(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Get requester ID from context
	requestorID, _ := ctx.Get("userID")

	progress, err := c.ProgressService.GetCourseStudentsProgress(uint(courseID), requestorID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// UpdateGradeByMentor handles updating a grade by a mentor
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

	// Get mentor ID from context
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
