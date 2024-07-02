package handlers

import (
	"net/http"
	"strconv"
	"time-tracker/client"
	"time-tracker/database"
	"time-tracker/models"

	"github.com/gin-gonic/gin"
)

// @Summary Add a new user
// @Description Add a new user with the given passport number
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user     body    models.User     true  "User"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Router /user [post]

// AddUser добавляет нового пользователя
func AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызов внешнего API для заполнения данных пользователя
	if err := client.FetchUserData(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data from external API"})
		return
	}

	// Сохранение в базу данных
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
		return
	}

	// Сохранение в базу данных
	database.DB.Create(&user)
	c.JSON(http.StatusOK, user)
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
// @Failure 400 {object} gin.H
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
