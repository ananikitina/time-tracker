package handlers

import (
	"fmt"
	"net/http"
	"time"
	"time-tracker/database"
	"time-tracker/models"

	"github.com/gin-gonic/gin"
)

func GetUserEffort(c *gin.Context) {
	// Extract parameters from the request
	userID := c.Param("userID")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, should be YYYY-MM-DD"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, should be YYYY-MM-DD"})
			return
		}
		// Adjust end date to include full day
		endDate = endDate.Add(24 * time.Hour).Add(-time.Second)
	}

	// Prepare query to calculate total effort (hours and minutes)
	var efforts []struct {
		UserID  string `json:"user_id"`
		Hours   int    `json:"hours"`
		Minutes int    `json:"minutes"`
	}

	query := database.DB.Model(&models.Task{}).
		Select("user_id, SUM(EXTRACT(HOUR FROM duration)) AS hours, SUM(EXTRACT(MINUTE FROM duration)) AS minutes").
		Where("user_id = ?", userID)

	// Apply date range filter if provided
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("start_time BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("start_time >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("start_time <= ?", endDate)
	}

	// Group by user_id and order by total effort (hours and minutes)
	query.Group("user_id").Order("hours DESC, minutes DESC").Scan(&efforts)

	// Log generated SQL query
	sql, _ := query.Debug().Find(&efforts).Rows()
	fmt.Println("SQL query:", sql)

	c.JSON(http.StatusOK, gin.H{
		"efforts": efforts,
	})
}

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
	now := time.Now().UTC()
	task.EndTime = &now

	// Save the updated task to the database
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

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
