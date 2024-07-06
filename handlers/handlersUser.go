package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ananikitina/time-tracker/database"
	"github.com/ananikitina/time-tracker/models"

	"github.com/gin-gonic/gin"
)

// error massage for swagger
type ErrorResponse struct {
	Error string `json:"error"`
}

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
// @Failure 400 {object} ErrorResponse "Invalid page parameter"
// @Failure 400 {object} ErrorResponse "Invalid pageSize parameter"
// @Failure 404 {object} ErrorResponse "No users found with specified filters"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	log.Println("Handling GetUsers request")

	var users []models.User
	query := database.DB

	// Filtering
	if passportNumber := c.Query("passportNumber"); passportNumber != "" {
		log.Printf("Filtering by passportNumber: %s", passportNumber)
		query = query.Where("passport_number = ?", passportNumber)
	}
	if surname := c.Query("surname"); surname != "" {
		log.Printf("Filtering by surname: %s", surname)
		query = query.Where("surname = ?", surname)
	}
	if name := c.Query("name"); name != "" {
		log.Printf("Filtering by name: %s", name)
		query = query.Where("name = ?", name)
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		log.Printf("Filtering by patronymic: %s", patronymic)
		query = query.Where("patronymic = ?", patronymic)
	}
	if address := c.Query("address"); address != "" {
		log.Printf("Filtering by address: %s", address)
		query = query.Where("address = ?", address)
	}

	// Pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("Error converting page to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		log.Printf("Error converting pageSize to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	offset := (page - 1) * pageSize
	log.Printf("Pagination - Page: %d, PageSize: %d, Offset: %d", page, pageSize, offset)

	query.Offset(offset).Limit(pageSize).Find(&users)
	log.Printf("Found %d users", len(users))

	// Check if users are found
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found with specified filters"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Add a new user
// @Description Add a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user     body    models.User     true  "User"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to save user to database"
// @Router /user [post]
func AddUser(c *gin.Context) {
	log.Println("Handling AddUser request")

	var newUser models.User

	// Parsing JSON request body
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Parsed user: %v", newUser)

	// Saving to database
	if err := database.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error saving user to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
		return
	}
	log.Printf("User saved: %v", newUser)

	c.JSON(http.StatusOK, newUser)
}

// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} ErrorResponse "User deleted successfully"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse  "Failed to delete user"
// @Router /user/{id} [delete]
func DeleteUser(c *gin.Context) {
	log.Println("Handling DeleteUser request")

	// Extracting user ID from URL parameter
	userID := c.Param("id")
	log.Printf("User ID to delete: %s", userID)

	// Checking if user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Deleting the user
	if err := database.DB.Delete(&user).Error; err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	log.Println("User deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// @Summary Update a user
// @Description Update a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body map[string]interface{} true "User data to update"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse "User not found"
// @Failure 404 {object} ErrorResponse "Invalid JSON format"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /user/{id} [put]
func UpdateUser(c *gin.Context) {
	log.Println("Handling UpdateUser request")

	var user models.User

	// Extracting user ID from URL parameter
	userID := c.Param("id")
	log.Printf("User ID to update: %s", userID)

	// Checking if user exists
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Binding JSON to the user structure
	var newUserData map[string]interface{}
	if err := c.ShouldBindJSON(&newUserData); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Printf("Updating user %s with data: %v", userID, newUserData)

	// Updating user's info
	if err := database.DB.Model(&user).Updates(newUserData).Error; err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	log.Println("User updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
