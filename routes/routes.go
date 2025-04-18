package routes

import (
	"LMS/controllers"
	middlewares "LMS/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// Protected routes (JWT middleware)
	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		// Courses
		protected.POST("/courses", controllers.CreateCourse)
		protected.GET("/courses", controllers.ListCourses)
		protected.GET("/courses/:id", controllers.GetCourse)

		// Materials
		protected.POST("/materials/upload", controllers.UploadMaterial)

		// Assignments
		protected.POST("/assignments", controllers.CreateAssignment)
		// (Endpoint untuk melihat tugas pada suatu course bisa ditambahkan sesuai kebutuhan)

		// Submissions
		protected.POST("/submissions", controllers.SubmitAssignment)

		// Enrollments
		protected.POST("/enrollments", controllers.EnrollCourse)
		protected.GET("/enrollments", controllers.GetMyEnrollments)

		// Assessments
		protected.POST("/assessments", controllers.CreateAssessment)

		// Discussions
		protected.POST("/discussions", controllers.CreateDiscussion)
		protected.GET("/discussions/course/:course_id", controllers.GetDiscussionsByCourse)
		protected.GET("/discussions/:id", controllers.GetDiscussionByID)

		// Comments
		protected.POST("/comments", controllers.CreateComment)
		protected.GET("/comments/discussion/:discussion_id", controllers.GetCommentsByDiscussion)

		}

	return router
}
