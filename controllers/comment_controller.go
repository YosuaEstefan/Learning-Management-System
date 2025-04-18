package controllers

import (
	"LMS/models"
	services "LMS/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment memungkinkan pengguna menambahkan komentar pada diskusi
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CommentServiceInstance.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByDiscussion menampilkan semua komentar untuk diskusi tertentu
func GetCommentsByDiscussion(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion id"})
		return
	}
	comments, err := services.CommentServiceInstance.GetCommentsByDiscussionID(uint(discussionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
