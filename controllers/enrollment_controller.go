package controllers

import (
	"LMS/models"
	services "LMS/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EnrollCourse memungkinkan siswa mendaftar ke kursus
func EnrollCourse(c *gin.Context) {
	if c.GetString("role") != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya siswa yang dapat mendaftar kursus"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	var input struct {
		CourseID uint `json:"course_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrollment := models.Enrollment{
		UserID:   userID.(uint),
		CourseID: input.CourseID,
	}
	if err := services.EnrollmentServiceInstance.Enroll(&enrollment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, enrollment)
}

// GetMyEnrollments menampilkan kursus yang telah didaftarkan siswa
func GetMyEnrollments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}
	enrollments, err := services.EnrollmentServiceInstance.GetEnrollmentsByUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}
