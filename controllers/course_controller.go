package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

)

// CourseController menangani permintaan kursus
type CourseController struct {
	CourseService *services.CourseService
}

// NewCourseController membuat pengontrol kursus baru
func NewCourseController(courseService *services.CourseService) *CourseController {
	return &CourseController{
		CourseService: courseService,
	}
}

// CreateCourseRequest mewakili permintaan untuk membuat kursus baru
type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	MentorID    uint   `json:"mentor_id" binding:"required"`
}

// UpdateCourseRequest mewakili permintaan untuk memperbarui kursus
type UpdateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// CreateCourse menangani pembuatan kursus
func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var request CreateCourseRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &models.Course{
		Title:       request.Title,
		Description: request.Description,
		MentorID:    request.MentorID,
	}

	if err := c.CourseService.CreateCourse(course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}

// GetCourseByID menangani pengambilan kursus berdasarkan ID
func (c *CourseController) GetCourseByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := c.CourseService.GetCourseByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	ctx.JSON(http.StatusOK, course)
}

// UpdateCourse menangani pembaruan kursus
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var request UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &models.Course{
		ID:          uint(id),
		Title:       request.Title,
		Description: request.Description,
	}

	if err := c.CourseService.UpdateCourse(course); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
	})
}

// DeleteCourse menangani penghapusan kursus
func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := c.CourseService.DeleteCourse(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course deleted successfully",
	})
}

// GetAllCourses menangani pengambilan semua kursus
func (c *CourseController) GetAllCourses(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	courses, err := c.CourseService.GetAllCourses(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

// GetCoursesByMentor menangani perolehan kursus berdasarkan ID mentor
func (c *CourseController) GetCoursesByMentor(ctx *gin.Context) {
	mentorID, err := strconv.ParseUint(ctx.Param("mentor_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mentor ID"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	courses, err := c.CourseService.GetCoursesByMentor(uint(mentorID), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}
