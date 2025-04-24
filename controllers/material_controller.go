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

// MaterialController menangani permintaan material
type MaterialController struct {
	MaterialService *services.MaterialService
}

// NewMaterialController membuat pengontrol material baru
func NewMaterialController(materialService *services.MaterialService) *MaterialController {
	return &MaterialController{
		MaterialService: materialService,
	}
}

// CreateMaterialRequest mewakili permintaan untuk membuat material baru
type CreateMaterialRequest struct {
	CourseID uint   `form:"course_id" binding:"required"`
	Title    string `form:"title" binding:"required"`
}

// UpdateMaterialRequest mewakili permintaan untuk memperbarui material
type UpdateMaterialRequest struct {
	Title string `form:"title" binding:"required"`
}

// CreateMaterial menangani pembuatan material
func (c *MaterialController) CreateMaterial(ctx *gin.Context) {
	var request CreateMaterialRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan file dari formulir
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Hasilkan nama file yang unik
	extension := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("material_%d_%s%s", request.CourseID, time.Now().Format("20060102150405"), extension)
	filePath := filepath.Join("uploads/materials", filename)

	// Simpan file ke disk
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

// GetMaterialByID menangani pengambilan material berdasarkan ID
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

// UpdateMaterial menangani pembaruan material
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
		// Hasilkan nama file yang unik
		extension := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("material_%d_%s%s", id, time.Now().Format("20060102150405"), extension)
		filePath = filepath.Join("uploads/materials", filename)

		// Simpan file ke disk
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

// DeleteMaterial menangani penghapusan material
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

// DownloadMaterial menangani pengunduhan material
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

	// Tetapkan nama file untuk diunduh
	filename := fmt.Sprintf("%s%s", material.Title, filepath.Ext(material.FilePath))
	ctx.FileAttachment(material.FilePath, filename)
}
