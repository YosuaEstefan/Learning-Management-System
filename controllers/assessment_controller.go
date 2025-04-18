package controllers

import (
	"LMS/models"
	services "LMS/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAssessment memungkinkan admin/mentor memberikan penilaian
func CreateAssessment(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "mentor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin atau mentor yang dapat memberikan penilaian"})
		return
	}

	var input models.Assessment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.AssessmentServiceInstance.CreateAssessment(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}
