package main

import (
	"golearn/config"
	"golearn/database"
	"golearn/handlers"
	"golearn/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "golearn/docs" // Swagger generated docs
)

// @title GoLearn API
// @version 1.0
// @description Uzaktan eğitim platformu backend API'si.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	// Initialize Database
	database.ConnectDB(cfg)

	router := gin.Default()

	// Global Middleware
	limiter := middleware.NewIPRateLimiter(5, cfg.RateLimitBurst)
	router.Use(middleware.RateLimitMiddleware(limiter))

	// WebSocket Hub
	hub := handlers.NewHub()
	go hub.Run()

	// Handlers
	authHandler := handlers.AuthHandler{Cfg: cfg}
	courseHandler := handlers.CourseHandler{}
	lessonHandler := handlers.LessonHandler{}
	quizHandler := handlers.QuizHandler{}
	progressHandler := handlers.ProgressHandler{}

	// Routes
	api := router.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Courses
			protected.GET("/courses", courseHandler.GetCourses)
			protected.GET("/courses/:id", courseHandler.GetCourse)
			protected.POST("/courses", middleware.TeacherOnly(), courseHandler.CreateCourse)
			protected.PUT("/courses/:id", middleware.TeacherOnly(), courseHandler.UpdateCourse)
			protected.DELETE("/courses/:id", middleware.TeacherOnly(), courseHandler.DeleteCourse)

			// Lessons
			protected.POST("/courses/:id/lessons", middleware.TeacherOnly(), lessonHandler.CreateLesson)
			protected.GET("/courses/:id/lessons", lessonHandler.GetLessons)

			// Quizzes
			protected.POST("/lessons/:id/quiz", middleware.TeacherOnly(), quizHandler.CreateQuiz)
			protected.GET("/lessons/:id/quiz", quizHandler.GetQuiz)
			protected.POST("/quiz/:id/submit", quizHandler.SubmitQuiz)

			// Progress
			protected.POST("/lessons/:id/complete", progressHandler.CompleteLesson)
			protected.GET("/my/progress", progressHandler.GetMyProgress)
		}
	}

	// WebSocket Route (Protected via context setting in handler or separate middleware)
	// We apply AuthMiddleware to ensure token is valid before connection
	router.GET("/ws/classroom/:courseId", middleware.AuthMiddleware(cfg), hub.HandleWebSocket)

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
