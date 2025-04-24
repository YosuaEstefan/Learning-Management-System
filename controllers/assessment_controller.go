package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AssessmentController handles assessment requests
type AssessmentController struct {
	AssessmentService *services.AssessmentService
}

// NewAssessmentController creates a new assessment controller
func NewAssessmentController(assessmentService *services.AssessmentService) *AssessmentController {
	return &AssessmentController{
		AssessmentService: assessmentService,
	}
}

// CreateAssessmentRequest represents a request to create a new assessment
type CreateAssessmentRequest struct {
	SubmissionID uint   `json:"submission_id" binding:"required"`
	Score        *int   `json:"score"`
	Feedback     string `json:"feedback"`
}

// UpdateAssessmentRequest represents a request to update an assessment
type UpdateAssessmentRequest struct {
	Score    *int   `json:"score"`
	Feedback string `json:"feedback"`
}

// CreateAssessment handles assessment creation
func (c *AssessmentController) CreateAssessment(ctx *gin.Context) {
	var request CreateAssessmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment := &models.Assessment{
		SubmissionID: request.SubmissionID,
		Score:        request.Score,
		Feedback:     request.Feedback,
	}

	if err := c.AssessmentService.CreateAssessment(assessment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Assessment created successfully",
		"assessment": assessment,
	})
}

// GetAssessmentByID handles getting an assessment by ID
func (c *AssessmentController) GetAssessmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assessment ID"})
		return
	}

	assessment, err := c.AssessmentService.GetAssessmentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	ctx.JSON(http.StatusOK, assessment)
}

// GetAssessmentBySubmission handles getting an assessment by submission ID
func (c *AssessmentController) GetAssessmentBySubmission(ctx *gin.Context) {
	submissionID, err := strconv.ParseUint(ctx.Param("submission_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	assessment, err := c.AssessmentService.GetAssessmentBySubmission(uint(submissionID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	ctx.JSON(http.StatusOK, assessment)
}

// UpdateAssessment handles updating an assessment
func (c *AssessmentController) UpdateAssessment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assessment ID"})
		return
	}

	var request UpdateAssessmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment := &models.Assessment{
		ID:       uint(id),
		Score:    request.Score,
		Feedback: request.Feedback,
	}

	if err := c.AssessmentService.UpdateAssessment(assessment); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Assessment updated successfully",
	})
}

// DeleteAssessment handles deleting an assessment
func (c *AssessmentController) DeleteAssessment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assessment ID"})
		return
	}

	if err := c.AssessmentService.DeleteAssessment(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Assessment deleted successfully",
	})
}
