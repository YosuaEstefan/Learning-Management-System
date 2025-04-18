package controllers

import (
	"LMS/config"
	"LMS/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAssignment memungkinkan mentor membuat tugas
func CreateAssignment(c *gin.Context) {
	// Hanya mentor yang boleh membuat tugas
	if c.GetString("role") != "mentor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya mentor yang dapat membuat tugas"})
		return
	}

	var assignment models.Assignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, assignment)
}
