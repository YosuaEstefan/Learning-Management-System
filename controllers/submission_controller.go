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

// SubmissionController menangani permintaan pengiriman
type SubmissionController struct {
	SubmissionService *services.SubmissionService
}

// NewSubmissionController membuat pengontrol pengiriman baru
func NewSubmissionController(submissionService *services.SubmissionService) *SubmissionController {
	return &SubmissionController{
		SubmissionService: submissionService,
	}
}

// CreateSubmissionRequest mewakili permintaan untuk membuat kiriman baru
type CreateSubmissionRequest struct {
	AssignmentID uint `form:"assignment_id" binding:"required"`
	StudentID    uint `form:"student_id" binding:"required"`
}

// CreateSubmission menangani pembuatan kiriman
func (c *SubmissionController) CreateSubmission(ctx *gin.Context) {
	var request CreateSubmissionRequest
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
	filename := fmt.Sprintf("submission_%d_%d_%s%s", request.AssignmentID, request.StudentID, time.Now().Format("20060102150405"), extension)
	filePath := filepath.Join("uploads/submissions", filename)

	// Menyimpan file ke disk
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	submission := &models.Submission{
		AssignmentID: request.AssignmentID,
		StudentID:    request.StudentID,
	}

	if err := c.SubmissionService.CreateSubmission(submission, filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Submission created successfully",
		"submission": submission,
	})
}

// GetSubmissionByID menangani mendapatkan kiriman dengan ID
func (c *SubmissionController) GetSubmissionByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := c.SubmissionService.GetSubmissionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	ctx.JSON(http.StatusOK, submission)
}

// GetSubmissionsByAssignment menangani mendapatkan kiriman berdasarkan ID penugasan
func (c *SubmissionController) GetSubmissionsByAssignment(ctx *gin.Context) {
	assignmentID, err := strconv.ParseUint(ctx.Param("assignment_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	submissions, err := c.SubmissionService.GetSubmissionsByAssignment(uint(assignmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get submissions"})
		return
	}

	ctx.JSON(http.StatusOK, submissions)
}

// GetSubmissionsByStudent menangani pengambilan kiriman berdasarkan ID siswa
func (c *SubmissionController) GetSubmissionsByStudent(ctx *gin.Context) {
	studentID, err := strconv.ParseUint(ctx.Param("student_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	submissions, err := c.SubmissionService.GetSubmissionsByStudent(uint(studentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get submissions"})
		return
	}

	ctx.JSON(http.StatusOK, submissions)
}

// DeleteSubmission menangani penghapusan kiriman
func (c *SubmissionController) DeleteSubmission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	if err := c.SubmissionService.DeleteSubmission(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Submission deleted successfully",
	})
}

// DownloadSubmission menangani pengunduhan kiriman
func (c *SubmissionController) DownloadSubmission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := c.SubmissionService.GetSubmissionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Tetapkan nama file untuk diunduh
	filename := fmt.Sprintf("submission_%d%s", submission.ID, filepath.Ext(submission.FilePath))
	ctx.FileAttachment(submission.FilePath, filename)
}
