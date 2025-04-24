// controllers/comment_controller.go
package controllers

import (
	"LMS/models"
	"LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentController handles comment requests
type CommentController struct {
	CommentService *services.CommentService
}

// NewCommentController creates a new comment controller
func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		CommentService: commentService,
	}
}

// CreateCommentRequest represents a request to create a new comment
type CreateCommentRequest struct {
	DiscussionID uint   `json:"discussion_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
}

// UpdateCommentRequest represents a request to update a comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment handles comment creation
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

// GetCommentByID handles getting a comment by ID
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

// GetCommentsByDiscussion handles getting comments by discussion ID
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

// UpdateComment handles updating a comment
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

// DeleteComment handles deleting a comment
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
