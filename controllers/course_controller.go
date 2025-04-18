package controllers

import (
	"LMS/config"
	"LMS/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCourse memungkinkan admin membuat kursus
func CreateCourse(c *gin.Context) {
	// Hanya admin yang boleh membuat kursus
	if c.GetString("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat membuat kursus"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

// ListCourses menampilkan semua kursus
func ListCourses(c *gin.Context) {
	var courses []models.Course
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetCourse menampilkan detail kursus berdasarkan id
func GetCourse(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var course models.Course
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kursus tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, course)
}
