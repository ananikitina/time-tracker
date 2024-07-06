package handlers

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/ananikitina/time-tracker/database"
	"github.com/ananikitina/time-tracker/models"

	"github.com/gin-gonic/gin"
)

// @Summary Sort user tasks
// @Description SortTasks sorts the user's tasks in descending order over a period of time.
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
	log.Println("Handling SortTasks request")

	var user models.User
	userID := c.Param("userID")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// Searching for a user by ID
	log.Printf("Finding user with ID: %s", userID)
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetching user's tasks for a specified time period
	log.Printf("Fetching tasks for user %s between %s and %s", userID, startTime, endTime)
	var tasks []models.Task
	if err := database.DB.Where("user_id = ? AND start_time >= ? AND end_time <= ?", user.ID, startTime, endTime).
		Find(&tasks).Error; err != nil {
		log.Printf("Failed to retrieve tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	// Calculating the duration of the task
	calculateDuration := func(task models.Task) time.Duration {
		if task.EndTime != nil {
			return task.EndTime.Sub(task.StartTime)
		}
		return time.Duration(0)
	}

	// Sorting tasks in descending order
	sortByDurationDesc := func(i, j int) bool {
		return calculateDuration(tasks[i]) > calculateDuration(tasks[j])
	}
	sort.Slice(tasks, sortByDurationDesc)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

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
	log.Println("Handling StartTask request")

	var user models.User
	userID := c.Param("userID")

	// Searching for a user by ID
	log.Printf("Finding user with ID: %s", userID)
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Creation of a new task
	task := models.Task{
		UserID:    user.ID,
		TaskName:  "Новая задача",
		StartTime: time.Now(),
		EndTime:   nil,
	}

	// Saving a task in a database
	log.Printf("Creating task for user %s: %+v", userID, task)
	if err := database.DB.Create(&task).Error; err != nil {
		log.Printf("Failed to create task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

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
	log.Println("Handling FinishTask request")

	var user models.User
	userID := c.Param("userID")

	// Searching for a user by ID
	log.Printf("Finding user with ID: %s", userID)
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Finding the active task for the user
	log.Printf("Finding active task for user %s", userID)
	var task models.Task
	if err := database.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&task).Error; err != nil {
		log.Printf("No active task found for user %s: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "No active task found for the user"})
		return
	}

	// Seting the end time to now
	now := time.Now()
	task.EndTime = &now

	// Saving the updated task to the database
	log.Printf("Finishing task for user %s: %+v", userID, task)
	if err := database.DB.Save(&task).Error; err != nil {
		log.Printf("Failed to finish task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

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
	log.Println("Handling GetUserTasks request")

	// Extract user ID from URL parameter
	userID := c.Param("id")
	log.Printf("Fetching tasks for user with ID: %s", userID)

	var tasks []models.Task

	// Fetching user's tasks
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Check if tasks are found
	if len(tasks) == 0 {
		log.Printf("No tasks found for user with ID: %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks found for the user"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
