package controllers

import (
	"LMS/models"
	services "LMS/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitAssignment memungkinkan siswa mengumpulkan tugas
func SubmitAssignment(c *gin.Context) {
	if c.GetString("role") != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya siswa yang dapat mengumpulkan tugas"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	var input struct {
		AssignmentID uint   `json:"assignment_id" binding:"required"`
		FilePath     string `json:"file_path" binding:"required"` // Simulasi file upload
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submission := models.Submission{
		AssignmentID: input.AssignmentID,
		StudentID:    userID.(uint),
		FilePath:     input.FilePath,
	}
	if err := services.SubmissionServiceInstance.CreateSubmission(&submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, submission)
}
