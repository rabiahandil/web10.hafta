package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct{}

// CreateLesson godoc
// @Summary Add a new lesson to a course
// @Tags lessons
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param body body models.Lesson true "Lesson details"
// @Success 201 {object} models.Lesson
// @Router /api/courses/{id}/lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	courseID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the course owner can add lessons"})
		return
	}

	var lesson models.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson.CourseID = course.ID
	if err := database.DB.Create(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create lesson"})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// GetLessons godoc
// @Summary List lessons of a course
// @Tags lessons
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {array} models.Lesson
// @Router /api/courses/{id}/lessons [get]
func (h *LessonHandler) GetLessons(c *gin.Context) {
	courseID := c.Param("id")
	var lessons []models.Lesson

	if err := database.DB.Where("course_id = ?", courseID).Order("\"order\" asc").Find(&lessons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch lessons"})
		return
	}

	c.JSON(http.StatusOK, lessons)
}
