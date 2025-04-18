package main

import (
	"LMS/config"
	"LMS/repositories"
	"LMS/routes"
	services "LMS/service"
)

func main() {
	// Koneksi database dan migrasi semua model
	db := config.ConnectDB()
	config.MigrateDB(db)

	// Inisialisasi service Enrollment
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	services.EnrollmentServiceInstance = services.NewEnrollmentService(enrollmentRepo)

	// Inisialisasi service Submission
	submissionRepo := repositories.NewSubmissionRepository(db)
	services.SubmissionServiceInstance = services.NewSubmissionService(submissionRepo)

	// Inisialisasi service Assessment
	assessmentRepo := repositories.NewAssessmentRepository(db)
	services.AssessmentServiceInstance = services.NewAssessmentService(assessmentRepo)

	// Inisialisasi service Discussion dan Comment
	discussionRepo := repositories.NewDiscussionRepository(db)
	services.DiscussionServiceInstance = services.NewDiscussionService(discussionRepo)

	commentRepo := repositories.NewCommentRepository(db)
	services.CommentServiceInstance = services.NewCommentService(commentRepo)

	router := routes.SetupRouter()
	router.Run(":8080")
}
