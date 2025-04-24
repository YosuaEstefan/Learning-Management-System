package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

)

// DiscussionController menangani permintaan diskusi
type DiscussionController struct {
	DiscussionService *services.DiscussionService
}

// NewDiscussionController membuat pengontrol diskusi baru
func NewDiscussionController(discussionService *services.DiscussionService) *DiscussionController {
	return &DiscussionController{
		DiscussionService: discussionService,
	}
}

// CreateDiscussionRequest mewakili permintaan untuk membuat diskusi baru
type CreateDiscussionRequest struct {
	CourseID uint   `json:"course_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// UpdateDiscussionRequest mewakili permintaan untuk memperbarui diskusi
type UpdateDiscussionRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreateDiscussion menangani pembuatan diskusi
func (c *DiscussionController) CreateDiscussion(ctx *gin.Context) {
	var request CreateDiscussionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	discussion := &models.Discussion{
		CourseID: request.CourseID,
		UserID:   userID.(uint),
		Title:    request.Title,
		Content:  request.Content,
	}

	if err := c.DiscussionService.CreateDiscussion(discussion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Discussion created successfully",
		"discussion": discussion,
	})
}

// GetDiscussionByID menangani mendapatkan diskusi berdasarkan ID
func (c *DiscussionController) GetDiscussionByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID"})
		return
	}

	discussion, err := c.DiscussionService.GetDiscussionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	ctx.JSON(http.StatusOK, discussion)
}

// GetDiscussionsByCourse menangani mendapatkan diskusi berdasarkan ID kursus
func (c *DiscussionController) GetDiscussionsByCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	discussions, err := c.DiscussionService.GetDiscussionsByCourse(uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discussions"})
		return
	}

	ctx.JSON(http.StatusOK, discussions)
}

// UpdateDiscussion menangani pembaruan diskusi
func (c *DiscussionController) UpdateDiscussion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID"})
		return
	}

	var request UpdateDiscussionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	discussion := &models.Discussion{
		ID:      uint(id),
		UserID:  userID.(uint),
		Title:   request.Title,
		Content: request.Content,
	}

	if err := c.DiscussionService.UpdateDiscussion(discussion); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Discussion updated successfully",
	})
}

// DeleteDiscussion menangani penghapusan diskusi
func (c *DiscussionController) DeleteDiscussion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID"})
		return
	}

	userID, _ := ctx.Get("userID")
	role, _ := ctx.Get("role")
	isAdmin := role == models.RoleAdmin

	if err := c.DiscussionService.DeleteDiscussion(uint(id), userID.(uint), isAdmin); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Discussion deleted successfully",
	})
}
