package handlers

import (
	"net/http"
	"sort"
	"time"
	"time-tracker/database"
	"time-tracker/models"

	"github.com/gin-gonic/gin"
)

// SortTasks сортирует задачи пользователя по убыванию продолжительности.
// @Summary Sort user tasks
// @Description Sort user tasks by duration in descending order
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param userID path string true "User ID"
// @Param start_time query string true "Start time filter (RFC3339 format)"
// @Param end_time query string true "End time filter (RFC3339 format)"
// @Success 200 {object} []models.Task
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve tasks"
// @Router /user/{userID}/tasks/sort [get]
func SortTasks(c *gin.Context) {
	var user models.User
	userID := c.Param("userID")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// Поиск пользователя по ID
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Извлечение задач пользователя за указанный период времени
	var tasks []models.Task
	if err := database.DB.Where("user_id = ? AND start_time >= ? AND end_time <= ?", user.ID, startTime, endTime).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	// Функция для расчета продолжительности задачи
	calculateDuration := func(task models.Task) time.Duration {
		if task.EndTime != nil {
			return task.EndTime.Sub(task.StartTime)
		}
		return time.Duration(0)
	}

	// Сортировка задач по убыванию продолжительности
	sortByDurationDesc := func(i, j int) bool {
		return calculateDuration(tasks[i]) > calculateDuration(tasks[j])
	}
	sort.Slice(tasks, sortByDurationDesc)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// StartTask начинает новую задачу для пользователя.
// @Summary Start a new task
// @Description Start a new task for the user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param userID path string true "User ID"
// @Success 201 {object} models.Task
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /user/{userID}/tasks/start [post]
func StartTask(c *gin.Context) {
	var user models.User
	userID := c.Param("userID")

	// Поиск пользователя по ID
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Создание новой задачи
	task := models.Task{
		UserID:    user.ID,
		TaskName:  "Новая задача",
		StartTime: time.Now(),
		EndTime:   nil, // Можно оставить nil или установить значение по умолчанию
	}

	// Сохранение задачи в базе данных
	database.DB.Create(&task)

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// FinishTask завершает активную задачу пользователя.
// @Summary Finish active task
// @Description Finish the active task for the user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param userID path string true "User ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 404 {object} ErrorResponse "No active task found for the user"
// @Failure 500 {object} ErrorResponse "Failed to finish task"
// @Router /user/{userID}/tasks/finish [put]
func FinishTask(c *gin.Context) {
	var user models.User
	userID := c.Param("userID")

	// Поиск пользователя по ID
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Find the active task for the user
	var task models.Task
	if err := database.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active task found for the user"})
		return
	}

	// Set the end time to now
	now := time.Now()
	task.EndTime = &now

	// Save the updated task to the database
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// GetUserTasks получает все задачи пользователя.
// @Summary Get user tasks
// @Description Get all tasks for the user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {array} models.Task
// @Failure 404 {object} ErrorResponse "No tasks found for the user"
// @Failure 500 {object} ErrorResponse "Failed to fetch tasks"
// @Router /user/{id}/tasks [get]
func GetUserTasks(c *gin.Context) {
	// Extract user ID from URL parameter
	userID := c.Param("id")

	var tasks []models.Task

	// Fetch all tasks for the user
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Check if tasks are found
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks found for the user"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
