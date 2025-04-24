package controllers

import (
	"LMS/models"
	services "LMS/services"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MaterialController handles material requests
type MaterialController struct {
	MaterialService *services.MaterialService
}

// NewMaterialController creates a new material controller
func NewMaterialController(materialService *services.MaterialService) *MaterialController {
	return &MaterialController{
		MaterialService: materialService,
	}
}

// CreateMaterialRequest represents a request to create a new material
type CreateMaterialRequest struct {
	CourseID uint   `form:"course_id" binding:"required"`
	Title    string `form:"title" binding:"required"`
}

// UpdateMaterialRequest represents a request to update a material
type UpdateMaterialRequest struct {
	Title string `form:"title" binding:"required"`
}

// CreateMaterial handles material creation
func (c *MaterialController) CreateMaterial(ctx *gin.Context) {
	var request CreateMaterialRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get file from form
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Generate unique filename
	extension := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("material_%d_%s%s", request.CourseID, time.Now().Format("20060102150405"), extension)
	filePath := filepath.Join("uploads/materials", filename)

	// Save file to disk
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	material := &models.Material{
		CourseID: request.CourseID,
		Title:    request.Title,
	}

	if err := c.MaterialService.CreateMaterial(material, filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Material created successfully",
		"material": material,
	})
}

// GetMaterialByID handles getting a material by ID
func (c *MaterialController) GetMaterialByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	material, err := c.MaterialService.GetMaterialByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	ctx.JSON(http.StatusOK, material)
}

// GetMaterialsByCourse handles getting materials by course ID
func (c *MaterialController) GetMaterialsByCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	materials, err := c.MaterialService.GetMaterialsByCourse(uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get materials"})
		return
	}

	ctx.JSON(http.StatusOK, materials)
}

// UpdateMaterial handles updating a material
func (c *MaterialController) UpdateMaterial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	var request UpdateMaterialRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material := &models.Material{
		ID:    uint(id),
		Title: request.Title,
	}

	filePath := ""
	file, err := ctx.FormFile("file")
	if err == nil {
		// Generate unique filename
		extension := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("material_%d_%s%s", id, time.Now().Format("20060102150405"), extension)
		filePath = filepath.Join("uploads/materials", filename)

		// Save file to disk
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	}

	if err := c.MaterialService.UpdateMaterial(material, filePath); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Material updated successfully",
	})
}

// DeleteMaterial handles deleting a material
func (c *MaterialController) DeleteMaterial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	if err := c.MaterialService.DeleteMaterial(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Material deleted successfully",
	})
}

// DownloadMaterial handles downloading a material
func (c *MaterialController) DownloadMaterial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	material, err := c.MaterialService.GetMaterialByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	// Set filename for download
	filename := fmt.Sprintf("%s%s", material.Title, filepath.Ext(material.FilePath))
	ctx.FileAttachment(material.FilePath, filename)
}
