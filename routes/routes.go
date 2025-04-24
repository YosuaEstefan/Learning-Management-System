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

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Create repositories
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

	// Create services
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

	// Create controllers
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

	// Create upload directories
	createUploadDirectories()

	// Routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", authController.GetProfile)

			// Courses
			courses := protected.Group("/courses")
			{
				// Routes for all authenticated users
				courses.GET("", courseController.GetAllCourses)
				courses.GET("/:id", courseController.GetCourseByID)

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				materials.GET("/:id", materialController.GetMaterialByID)
				materials.GET("/course/:course_id", materialController.GetMaterialsByCourse)
				materials.GET("/download/:id", materialController.DownloadMaterial)

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				assignments.GET("/:id", assignmentController.GetAssignmentByID)
				assignments.GET("/course/:course_id", assignmentController.GetAssignmentsByCourse)

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				enrollments.GET("/:id", enrollmentController.GetEnrollmentByID)
				enrollments.GET("/check", enrollmentController.CheckEnrollment)

				// Routes for students
				studentEnrollments := enrollments.Group("/")
				studentEnrollments.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleStudent)(c)
				})
				{
					studentEnrollments.POST("", enrollmentController.CreateEnrollment)
					studentEnrollments.DELETE("/:id", enrollmentController.DeleteEnrollment)
				}

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				submissions.GET("/:id", submissionController.GetSubmissionByID)
				submissions.GET("/download/:id", submissionController.DownloadSubmission)

				// Routes for students
				studentSubmissions := submissions.Group("/")
				studentSubmissions.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleStudent)(c)
				})
				{
					studentSubmissions.POST("", submissionController.CreateSubmission)
					studentSubmissions.DELETE("/:id", submissionController.DeleteSubmission)
				}

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				assessments.GET("/:id", assessmentController.GetAssessmentByID)
				assessments.GET("/submission/:submission_id", assessmentController.GetAssessmentBySubmission)

				// Routes for admins and mentors
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
				// Routes for all authenticated users
				discussions.GET("/:id", discussionController.GetDiscussionByID)
				discussions.GET("/course/:course_id", discussionController.GetDiscussionsByCourse)
				discussions.POST("", discussionController.CreateDiscussion)
				discussions.PUT("/:id", discussionController.UpdateDiscussion)
				discussions.DELETE("/:id", discussionController.DeleteDiscussion)
			}

			// Comments
			comments := protected.Group("/comments")
			{
				// Routes for all authenticated users
				comments.GET("/:id", commentController.GetCommentByID)
				comments.GET("/discussion/:discussion_id", commentController.GetCommentsByDiscussion)
				comments.POST("", commentController.CreateComment)
				comments.PUT("/:id", commentController.UpdateComment)
				comments.DELETE("/:id", commentController.DeleteComment)
			}
			// Learning Progress routes
			progress := protected.Group("/progress")
			{
				// Routes for all authenticated users
				progress.GET("/:id", progressController.GetProgressByID)
				progress.GET("/student/:student_id/course/:course_id", progressController.GetStudentProgress)

				// Routes for mentors and admins
				mentorAdminProgress := progress.Group("/")
				mentorAdminProgress.Use(func(c *gin.Context) {
					middleware.RoleMiddleware(models.RoleAdmin, models.RoleMentor)(c)
				})
				{
					mentorAdminProgress.POST("/grade", progressController.GradeActivity)
					mentorAdminProgress.GET("/course/:course_id/students", progressController.GetCourseStudentsProgress)
				}

				// Routes only for mentors (updating grades)
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

// Helper function to create upload directories
func createUploadDirectories() {
	// Create upload directories if they don't exist
	// This would be implemented with os.MkdirAll to create directories
	// for storing uploaded files
}
