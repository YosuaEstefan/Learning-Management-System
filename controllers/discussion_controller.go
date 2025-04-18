package controllers

import (
	"LMS/models"
	services "LMS/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDiscussion memungkinkan admin, mentor, atau student membuat diskusi
func CreateDiscussion(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "mentor" && role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak diizinkan membuat diskusi"})
		return
	}

	var discussion models.Discussion
	if err := c.ShouldBindJSON(&discussion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.DiscussionServiceInstance.CreateDiscussion(&discussion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, discussion)
}

// GetDiscussionByID mengambil detail diskusi berdasarkan id
func GetDiscussionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion id"})
		return
	}
	discussion, err := services.DiscussionServiceInstance.GetDiscussionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diskusi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, discussion)
}

// GetDiscussionsByCourse menampilkan daftar diskusi berdasarkan course_id
func GetDiscussionsByCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	discussions, err := services.DiscussionServiceInstance.GetDiscussionsByCourseID(uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, discussions)
}
