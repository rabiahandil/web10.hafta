package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct{}

// CreateCourse godoc
// @Summary Create a new course
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Course true "Course details"
// @Success 201 {object} models.Course
// @Router /api/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	course.TeacherID = userID

	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourses godoc
// @Summary List courses with pagination and filtering
// @Tags courses
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param category query string false "Filter by category"
// @Param sort query string false "Sort by field (e.g. title, created_at)"
// @Success 200 {object} map[string]interface{}
// @Router /api/courses [get]
func (h *CourseHandler) GetCourses(c *gin.Context) {
	var courses []models.Course
	var total int64

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}
	
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at desc")

	query := database.DB.Model(&models.Course{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not count courses"})
		return
	}

	offset := (page - 1) * limit
	if err := query.Order(sort).Limit(limit).Offset(offset).Preload("Teacher").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetCourse godoc
// @Summary Get course details
// @Tags courses
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Router /api/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := database.DB.Preload("Teacher").Preload("Lessons").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// UpdateCourse godoc
// @Summary Update an existing course
// @Tags courses
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param body body models.Course true "Update details"
// @Success 200 {object} models.Course
// @Router /api/courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own courses"})
		return
	}

	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&course).Updates(input)
	c.JSON(http.StatusOK, course)
}

// DeleteCourse godoc
// @Summary Delete a course
// @Tags courses
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]string
// @Router /api/courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own courses"})
		return
	}

	database.DB.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
