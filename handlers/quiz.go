package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct{}

// CreateQuiz godoc
// @Summary Create a quiz for a lesson
// @Tags quizzes
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Param body body models.Quiz true "Quiz details"
// @Success 21 {object} models.Quiz
// @Router /api/lessons/{id}/quiz [post]
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	lessonID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	// Check if user is the teacher of the course
	var course models.Course
	if err := database.DB.First(&course, lesson.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if course.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the course owner can create quizzes"})
		return
	}

	var quiz models.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz.LessonID = lesson.ID
	if err := database.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create quiz"})
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

// GetQuiz godoc
// @Summary Get quiz details by lesson ID
// @Tags quizzes
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.Quiz
// @Router /api/lessons/{id}/quiz [get]
func (h *QuizHandler) GetQuiz(c *gin.Context) {
	lessonID := c.Param("id")
	var quiz models.Quiz

	if err := database.DB.Preload("Questions").Where("lesson_id = ?", lessonID).First(&quiz).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

type QuizSubmission struct {
	Answers map[string]string `json:"answers"` // question_id -> answer
}

// SubmitQuiz godoc
// @Summary Submit quiz answers and get score
// @Tags quizzes
// @Security BearerAuth
// @Param id path int true "Quiz ID"
// @Param body body QuizSubmission true "Answers map"
// @Success 200 {object} map[string]interface{}
// @Router /api/quiz/{id}/submit [post]
func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	quizID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var quiz models.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	var submission QuizSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	correctCount := 0
	totalQuestions := len(quiz.Questions)

	for _, q := range quiz.Questions {
		qIDStr := strconv.FormatUint(uint64(q.ID), 10)
		if ans, ok := submission.Answers[qIDStr]; ok && ans == q.Correct {
			correctCount++
		}
	}

	percent := 0.0
	if totalQuestions > 0 {
		percent = (float64(correctCount) / float64(totalQuestions)) * 100
	}

	result := models.QuizResult{
		UserID:  userID,
		QuizID:  quiz.ID,
		Score:   correctCount,
		Total:   totalQuestions,
		Percent: percent,
	}

	database.DB.Create(&result)

	c.JSON(http.StatusOK, gin.H{
		"score":   correctCount,
		"total":   totalQuestions,
		"percent": percent,
		"message": "Quiz submitted successfully",
	})
}
