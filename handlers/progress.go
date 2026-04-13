package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct{}

// CompleteLesson godoc
// @Summary Mark a lesson as completed
// @Tags progress
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/lessons/{id}/complete [post]
func (h *ProgressHandler) CompleteLesson(c *gin.Context) {
	lessonID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	progress := models.Progress{
		UserID:   userID,
		LessonID: lesson.ID,
		CourseID: lesson.CourseID,
	}

	// Use Create or update (duplicate handled by unique index in DB, but we check here for cleaner response)
	var existing models.Progress
	if err := database.DB.Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Lesson already completed"})
		return
	}

	if err := database.DB.Create(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not mark lesson as complete"})
		return
	}

	// Calculate and return current course progress
	h.GetMyProgress(c) // Re-use the progress listing logic for immediate feedback
}

// GetMyProgress godoc
// @Summary Get progress of all enrolled courses
// @Tags progress
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /api/my/progress [get]
func (h *ProgressHandler) GetMyProgress(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	
	type ProgressResult struct {
		CourseID         uint    `json:"course_id"`
		CourseTitle      string  `json:"course_title"`
		TotalLessons     int64   `json:"total_lessons"`
		CompletedLessons int64   `json:"completed_lessons"`
		Percent          float64 `json:"percent"`
	}

	var results []ProgressResult

	// Get all courses that the user has some progress in
	var courseIDs []uint
	database.DB.Model(&models.Progress{}).Where("user_id = ?", userID).Distinct("course_id").Pluck("course_id", &courseIDs)

	for _, cid := range courseIDs {
		var course models.Course
		database.DB.First(&course, cid)

		var totalLessons int64
		database.DB.Model(&models.Lesson{}).Where("course_id = ?", cid).Count(&totalLessons)

		var completedLessons int64
		database.DB.Model(&models.Progress{}).Where("user_id = ? AND course_id = ?", userID, cid).Count(&completedLessons)

		percent := 0.0
		if totalLessons > 0 {
			percent = (float64(completedLessons) / float64(totalLessons)) * 100
		}

		results = append(results, ProgressResult{
			CourseID:         cid,
			CourseTitle:      course.Title,
			TotalLessons:     totalLessons,
			CompletedLessons: completedLessons,
			Percent:          percent,
		})
	}

	c.JSON(http.StatusOK, results)
}
