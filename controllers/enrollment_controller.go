package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

)

// EnrollmentController menangani permintaan pendaftaran
type EnrollmentController struct {
	EnrollmentService *services.EnrollmentService
}

// NewEnrollmentController membuat pengontrol pendaftaran baru
func NewEnrollmentController(enrollmentService *services.EnrollmentService) *EnrollmentController {
	return &EnrollmentController{
		EnrollmentService: enrollmentService,
	}
}

// CreateEnrollmentRequest mewakili permintaan untuk membuat pendaftaran baru
type CreateEnrollmentRequest struct {
	UserID   uint `json:"user_id" binding:"required"`
	CourseID uint `json:"course_id" binding:"required"`
}

// CreateEnrollment menangani pembuatan pendaftaran
func (c *EnrollmentController) CreateEnrollment(ctx *gin.Context) {
	var request CreateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrollment := &models.Enrollment{
		UserID:   request.UserID,
		CourseID: request.CourseID,
	}

	if err := c.EnrollmentService.CreateEnrollment(enrollment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Enrollment created successfully",
		"enrollment": enrollment,
	})
}

// GetEnrollmentByID menangani pendaftaran berdasarkan ID
func (c *EnrollmentController) GetEnrollmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	enrollment, err := c.EnrollmentService.GetEnrollmentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	ctx.JSON(http.StatusOK, enrollment)
}

// GetEnrollmentsByUser menangani pendaftaran berdasarkan ID pengguna
func (c *EnrollmentController) GetEnrollmentsByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	enrollments, err := c.EnrollmentService.GetEnrollmentsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
		return
	}

	ctx.JSON(http.StatusOK, enrollments)
}

// GetEnrollmentsByCourse menangani pendaftaran berdasarkan ID kursus
func (c *EnrollmentController) GetEnrollmentsByCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	enrollments, err := c.EnrollmentService.GetEnrollmentsByCourse(uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
		return
	}

	ctx.JSON(http.StatusOK, enrollments)
}

// DeleteEnrollment menangani penghapusan pendaftaran
func (c *EnrollmentController) DeleteEnrollment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	if err := c.EnrollmentService.DeleteEnrollment(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Enrollment deleted successfully",
	})
}

// CheckEnrollment menangani pengecekan apakah seorang siswa terdaftar dalam suatu kursus
func (c *EnrollmentController) CheckEnrollment(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Query("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	courseID, err := strconv.ParseUint(ctx.Query("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	enrolled, err := c.EnrollmentService.IsStudentEnrolled(uint(userID), uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"enrolled": enrolled,
	})
}
