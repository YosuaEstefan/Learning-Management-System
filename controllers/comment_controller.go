package controllers

import (
	"LMS/models"
	"LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

)

// CommentController menangani permintaan komentar
type CommentController struct {
	CommentService *services.CommentService
}

// NewCommentController membuat pengontrol komentar baru
func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		CommentService: commentService,
	}
}

// CreateCommentRequest mewakili permintaan untuk membuat komentar baru
type CreateCommentRequest struct {
	DiscussionID uint   `json:"discussion_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
}

// UpdateCommentRequest mewakili permintaan untuk memperbarui komentar
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment menangani pembuatan komentar
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var request CreateCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	comment := &models.Comment{
		DiscussionID: request.DiscussionID,
		UserID:       userID.(uint),
		Content:      request.Content,
	}

	if err := c.CommentService.CreateComment(comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetCommentByID menangani pengambilan komentar berdasarkan ID
func (c *CommentController) GetCommentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := c.CommentService.GetCommentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// GetCommentsByDiscussion menangani pengambilan komentar berdasarkan ID diskusi
func (c *CommentController) GetCommentsByDiscussion(ctx *gin.Context) {
	discussionID, err := strconv.ParseUint(ctx.Param("discussion_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID"})
		return
	}

	comments, err := c.CommentService.GetCommentsByDiscussion(uint(discussionID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// UpdateComment menangani pembaruan komentar
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var request UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	comment := &models.Comment{
		ID:      uint(id),
		UserID:  userID.(uint),
		Content: request.Content,
	}

	if err := c.CommentService.UpdateComment(comment); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
	})
}

// DeleteComment menangani penghapusan komentar
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, _ := ctx.Get("userID")
	role, _ := ctx.Get("role")
	isAdmin := role == models.RoleAdmin

	if err := c.CommentService.DeleteComment(uint(id), userID.(uint), isAdmin); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
