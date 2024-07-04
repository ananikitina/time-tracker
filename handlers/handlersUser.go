package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time-tracker/database"
	"time-tracker/models"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get users
// @Description Get users with filtering and pagination
// @Tags users
// @Accept  json
// @Produce  json
// @Param passportNumber query string false "Passport Number"
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {array} models.User
// @Failure 400 {object} map[string]interface{}
// @Router /users [get]

// GetUsers получает данные пользователя
func GetUsers(c *gin.Context) {
	var users []models.User
	query := database.DB

	// Фильтрация по всем полям
	if passportNumber := c.Query("passportNumber"); passportNumber != "" {
		query = query.Where("passport_number = ?", passportNumber)
	}
	if surname := c.Query("surname"); surname != "" {
		query = query.Where("surname = ?", surname)
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name = ?", name)
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		query = query.Where("patronymic = ?", patronymic)
	}
	if address := c.Query("address"); address != "" {
		query = query.Where("address = ?", address)
	}

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10")) //количество элементов на одной странице
	offset := (page - 1) * pageSize

	query.Offset(offset).Limit(pageSize).Find(&users)
	c.JSON(http.StatusOK, users)
}

// AddUser godoc
// @Summary Add a new user
// @Description Add a new user with the given passport number
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user     body    models.User     true  "User"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user [post]
func AddUser(c *gin.Context) {
	var newUser models.User

	// Parse JSON request body
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Save to database
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
		return
	}

	// Return the created user as JSON response
	c.JSON(http.StatusOK, newUser)
}

func DeleteUser(c *gin.Context) {
	// Extract user ID from URL parameter
	userID := c.Param("id")

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func UpdateUser(c *gin.Context) {
	var user models.User

	// Получить ID пользователя из параметров URL
	userID := c.Param("id")

	// Проверить существование пользователя в базе данных
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Принять данные JSON из запроса и привязать их к структуре User
	var newUserData map[string]interface{}
	if err := c.ShouldBindJSON(&newUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Логирование для отладки
	log.Printf("Updating user %s with data: %v", userID, newUserData)

	// Обновить данные пользователя, используя карту
	if err := database.DB.Model(&user).Updates(newUserData).Error; err != nil {
		log.Printf("Error updating user: %v", err) // Логирование ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
