package routes

import (
	"LMS/controllers"
	"LMS/middleware"
	"LMS/models"
	"LMS/repositories"
	"LMS/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)

// SetupRoutes mengkonfigurasi rute API
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Buat repositori
	userRepo := repositories.NewUserRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	materialRepo := repositories.NewMaterialRepository(db)
	assignmentRepo := repositories.NewAssignmentRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)
	assessmentRepo := repositories.NewAssessmentRepository(db)
	discussionRepo := repositories.NewDiscussionRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	progressRepo := repositories.NewLearningProgressRepository(db)

	// buat service
	authService := services.NewAuthService(userRepo)
	courseService := services.NewCourseService(courseRepo, userRepo)
	materialService := services.NewMaterialService(materialRepo, courseRepo)
	assignmentService := services.NewAssignmentService(assignmentRepo, courseRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, userRepo, courseRepo)
	submissionService := services.NewSubmissionService(submissionRepo, assignmentRepo, enrollmentRepo, userRepo)
	assessmentService := services.NewAssessmentService(assessmentRepo, submissionRepo, assignmentRepo, userRepo)
	discussionService := services.NewDiscussionService(discussionRepo, courseRepo, userRepo, enrollmentRepo)
	commentService := services.NewCommentService(commentRepo, discussionRepo, userRepo, courseRepo, enrollmentRepo)
	progressService := services.NewLearningProgressService(progressRepo, userRepo, courseRepo, enrollmentRepo, assignmentRepo)

	// buat controllers
	authController := controllers.NewAuthController(authService)
	courseController := controllers.NewCourseController(courseService)
	materialController := controllers.NewMaterialController(materialService)
	assignmentController := controllers.NewAssignmentController(assignmentService)
	enrollmentController := controllers.NewEnrollmentController(enrollmentService)
	submissionController := controllers.NewSubmissionController(submissionService)
	assessmentController := controllers.NewAssessmentController(assessmentService)
	discussionController := controllers.NewDiscussionController(discussionService)
	commentController := controllers.NewCommentController(commentService)
	progressController := controllers.NewProgressController(progressService)

	// Buat direktori unggahan
	createUploadDirectories()

	// router
	api := router.Group("/api")
	{
		// Publik router
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Rute yang dilindungi
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profil pengguna
			protected.GET("/profile", authController.GetProfile)

			// Courses
			courses := protected.Group("/courses")
			{
				//Rute untuk para admin dan mentor
				courses.GET("", courseController.GetAllCourses)
				courses.GET("/:id", courseController.GetCourseByID)

				//Rute untuk para admin dan mentor
				adminMentorCourses := courses.Group("/")
				adminMentorCourses.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorCourses.POST("", courseController.CreateCourse)
					adminMentorCourses.PUT("/:id", courseController.UpdateCourse)
					adminMentorCourses.DELETE("/:id", courseController.DeleteCourse)
					adminMentorCourses.GET("/mentor/:mentor_id", courseController.GetCoursesByMentor)
				}
			}

			// Materials
			materials := protected.Group("/materials")
			{
				// Rute untuk semua pengguna yang diautentikasi
				materials.GET("/:id", materialController.GetMaterialByID)
				materials.GET("/course/:course_id", materialController.GetMaterialsByCourse)
				materials.GET("/download/:id", materialController.DownloadMaterial)
				// Rute untuk admin dan mentor
				adminMentorMaterials := materials.Group("/")
				adminMentorMaterials.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorMaterials.POST("", materialController.CreateMaterial)
					adminMentorMaterials.PUT("/:id", materialController.UpdateMaterial)
					adminMentorMaterials.DELETE("/:id", materialController.DeleteMaterial)
				}
			}

			// Assignments
			assignments := protected.Group("/assignments")
			{
				// Rute untuk semua pengguna yang diautentikasi
				assignments.GET("/:id", assignmentController.GetAssignmentByID)
				assignments.GET("/course/:course_id", assignmentController.GetAssignmentsByCourse)

				// Rute untuk admin dan mentor
				adminMentorAssignments := assignments.Group("/")
				adminMentorAssignments.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorAssignments.POST("", assignmentController.CreateAssignment)
					adminMentorAssignments.PUT("/:id", assignmentController.UpdateAssignment)
					adminMentorAssignments.DELETE("/:id", assignmentController.DeleteAssignment)
				}
			}

			// Enrollments
			enrollments := protected.Group("/enrollments")
			{
				// Rute untuk semua pengguna yang diautentikasi
				enrollments.GET("/:id", enrollmentController.GetEnrollmentByID)
				enrollments.GET("/check", enrollmentController.CheckEnrollment)

				// Rute untuk siswa
				studentEnrollments := enrollments.Group("/")
				studentEnrollments.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleStudent)(c)
				})
				{
					studentEnrollments.POST("", enrollmentController.CreateEnrollment)
					studentEnrollments.DELETE("/:id", enrollmentController.DeleteEnrollment)
				}

				// Rute untuk admin dan mentor
				adminMentorEnrollments := enrollments.Group("/")
				adminMentorEnrollments.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorEnrollments.GET("/user/:user_id", enrollmentController.GetEnrollmentsByUser)
					adminMentorEnrollments.GET("/course/:course_id", enrollmentController.GetEnrollmentsByCourse)
				}
			}

			// Submissions
			submissions := protected.Group("/submissions")
			{
				// Rute untuk semua pengguna yang diautentikasi
				submissions.GET("/:id", submissionController.GetSubmissionByID)
				submissions.GET("/download/:id", submissionController.DownloadSubmission)

				// Rute untuk siswa
				studentSubmissions := submissions.Group("/")
				studentSubmissions.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleStudent)(c)
				})
				{
					studentSubmissions.POST("", submissionController.CreateSubmission)
					studentSubmissions.DELETE("/:id", submissionController.DeleteSubmission)
				}

				// Rute untuk admin dan mentor
				adminMentorSubmissions := submissions.Group("/")
				adminMentorSubmissions.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorSubmissions.GET("/assignment/:assignment_id", submissionController.GetSubmissionsByAssignment)
					adminMentorSubmissions.GET("/student/:student_id", submissionController.GetSubmissionsByStudent)
				}
			}

			// Assessments
			assessments := protected.Group("/assessments")
			{
				// Rute untuk semua pengguna yang diautentikasi
				assessments.GET("/:id", assessmentController.GetAssessmentByID)
				assessments.GET("/submission/:submission_id", assessmentController.GetAssessmentBySubmission)

				// Rute untuk admin dan mentor
				adminMentorAssessments := assessments.Group("/")
				adminMentorAssessments.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					adminMentorAssessments.POST("", assessmentController.CreateAssessment)
					adminMentorAssessments.PUT("/:id", assessmentController.UpdateAssessment)
					adminMentorAssessments.DELETE("/:id", assessmentController.DeleteAssessment)
				}
			}

			// Discussions
			discussions := protected.Group("/discussions")
			{
				// Rute untuk semua pengguna yang diautentikasi
				discussions.GET("/:id", discussionController.GetDiscussionByID)
				discussions.GET("/course/:course_id", discussionController.GetDiscussionsByCourse)
				discussions.POST("", discussionController.CreateDiscussion)
				discussions.PUT("/:id", discussionController.UpdateDiscussion)
				discussions.DELETE("/:id", discussionController.DeleteDiscussion)
			}

			// Comments
			comments := protected.Group("/comments")
			{
				// Rute untuk semua pengguna yang diautentikasi
				comments.GET("/:id", commentController.GetCommentByID)
				comments.GET("/discussion/:discussion_id", commentController.GetCommentsByDiscussion)
				comments.POST("", commentController.CreateComment)
				comments.PUT("/:id", commentController.UpdateComment)
				comments.DELETE("/:id", commentController.DeleteComment)
			}
			// Learning Progress
			progress := protected.Group("/progress")
			{
				// Rute untuk semua pengguna yang diautentikasi
				progress.GET("/:id", progressController.GetProgressByID)
				progress.GET("/student/:student_id/course/:course_id", progressController.GetStudentProgress)

				// Rute untuk mentor dan admin
				mentorAdminProgress := progress.Group("/")
				mentorAdminProgress.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					mentorAdminProgress.POST("/grade", progressController.GradeActivity)
					mentorAdminProgress.GET("/course/:course_id/students", progressController.GetCourseStudentsProgress)
				}

				// Rute hanya untuk mentor
				mentorProgress := progress.Group("/")
				mentorProgress.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleMentor)(c)
				})
				{
					mentorProgress.PUT("/update-grade/:progress_id", progressController.UpdateGradeByMentor)
				}
			}

		}
	}
}

// Fungsi pembantu untuk membuat direktori unggahan
func createUploadDirectories() {

}
