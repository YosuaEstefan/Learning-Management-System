package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EnrollmentController handles enrollment requests
type EnrollmentController struct {
	EnrollmentService *services.EnrollmentService
}

// NewEnrollmentController creates a new enrollment controller
func NewEnrollmentController(enrollmentService *services.EnrollmentService) *EnrollmentController {
	return &EnrollmentController{
		EnrollmentService: enrollmentService,
	}
}

// CreateEnrollmentRequest represents a request to create a new enrollment
type CreateEnrollmentRequest struct {
	UserID   uint `json:"user_id" binding:"required"`
	CourseID uint `json:"course_id" binding:"required"`
}

// CreateEnrollment handles enrollment creation
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

// GetEnrollmentByID handles getting an enrollment by ID
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

// GetEnrollmentsByUser handles getting enrollments by user ID
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

// GetEnrollmentsByCourse handles getting enrollments by course ID
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

// DeleteEnrollment handles deleting an enrollment
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

// CheckEnrollment handles checking if a student is enrolled in a course
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
