package controllers

import (
	"LMS/models"
	services "LMS/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AssignmentController handles assignment requests
type AssignmentController struct {
	AssignmentService *services.AssignmentService
}

// NewAssignmentController creates a new assignment controller
func NewAssignmentController(assignmentService *services.AssignmentService) *AssignmentController {
	return &AssignmentController{
		AssignmentService: assignmentService,
	}
}

// CreateAssignmentRequest represents a request to create a new assignment
type CreateAssignmentRequest struct {
	CourseID    uint       `json:"course_id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	MaxScore    *int       `json:"max_score"`
}

// UpdateAssignmentRequest represents a request to update an assignment
type UpdateAssignmentRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	MaxScore    *int       `json:"max_score"`
}

// CreateAssignment handles assignment creation
func (c *AssignmentController) CreateAssignment(ctx *gin.Context) {
	var request CreateAssignmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment := &models.Assignment{
		CourseID:    request.CourseID,
		Title:       request.Title,
		Description: request.Description,
		DueDate:     request.DueDate,
		MaxScore:    request.MaxScore,
	}

	if err := c.AssignmentService.CreateAssignment(assignment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Assignment created successfully",
		"assignment": assignment,
	})
}

// GetAssignmentByID handles getting an assignment by ID
func (c *AssignmentController) GetAssignmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := c.AssignmentService.GetAssignmentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	ctx.JSON(http.StatusOK, assignment)
}

// GetAssignmentsByCourse handles getting assignments by course ID
func (c *AssignmentController) GetAssignmentsByCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("course_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	assignments, err := c.AssignmentService.GetAssignmentsByCourse(uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get assignments"})
		return
	}

	ctx.JSON(http.StatusOK, assignments)
}

// UpdateAssignment handles updating an assignment
func (c *AssignmentController) UpdateAssignment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var request UpdateAssignmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment := &models.Assignment{
		ID:          uint(id),
		Title:       request.Title,
		Description: request.Description,
		DueDate:     request.DueDate,
		MaxScore:    request.MaxScore,
	}

	if err := c.AssignmentService.UpdateAssignment(assignment); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Assignment updated successfully",
	})
}

// DeleteAssignment handles deleting an assignment
func (c *AssignmentController) DeleteAssignment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := c.AssignmentService.DeleteAssignment(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Assignment deleted successfully",
	})
}
