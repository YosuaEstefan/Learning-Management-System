package controllers

import (
	"LMS/config"
	"LMS/models"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadMaterial oleh mentor untuk mengunggah materi
func UploadMaterial(c *gin.Context) {
	// Hanya mentor yang dapat mengunggah materi
	if c.GetString("role") != "mentor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya mentor yang dapat mengunggah materi"})
		return
	}

	courseIDStr := c.PostForm("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
		return
	}
	title := c.PostForm("title")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saat mengunggah file: " + err.Error()})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".mp4" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe file tidak diperbolehkan"})
		return
	}

	filePath := "uploads/materials/" + time.Now().Format("20060102150405") + "_" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file: " + err.Error()})
		return
	}

	material := models.Material{
		CourseID: uint(courseID),
		Title:    title,
		FilePath: filePath,
	}
	if err := config.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, material)
}
