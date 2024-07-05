package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"time-tracker/database"
	"time-tracker/models"
	"time-tracker/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestUser возвращает тестового пользователя
func getTestUser() models.User {
	return models.User{
		PassportNumber: "890231",
		Surname:        "Вавилов",
		Name:           "Анатолий",
		Patronymic:     "Анатольевич",
		Address:        "г.Москва, ул. Кирова д.19",
	}
}
func setupRouter() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	database.Connect()
	database.Migrate()

	// Инициализация Gin
	r := gin.Default()

	// Регистрация маршрутов
	routes.SetupRouter(r)

	return r
}

// TestAddGetDelete проверяет добавление и удаление пользователя
func TestAddGetDelete(t *testing.T) {
	// prepare
	router := setupRouter()

	// Создание тестового пользователя
	user := getTestUser()

	// Создание тестового HTTP-запроса
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Запись HTTP-ответа
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, w.Code)

	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err, "failed to unmarshal JSON response")

	// Проверка совпадения созданного пользователя с ожидаемыми данными
	assert.Equal(t, user.PassportNumber, createdUser.PassportNumber)
	assert.Equal(t, user.Surname, createdUser.Surname)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Patronymic, createdUser.Patronymic)
	assert.Equal(t, user.Address, createdUser.Address)

	// Проверка, что пользователь был добавлен в базу данных
	var dbUser models.User
	err = database.DB.Where("passport_number = ?", user.PassportNumber).First(&dbUser).Error
	require.NoError(t, err, "failed to fetch user from database")
	assert.Equal(t, user.PassportNumber, dbUser.PassportNumber)
	assert.Equal(t, user.Surname, dbUser.Surname)
	assert.Equal(t, user.Name, dbUser.Name)
	assert.Equal(t, user.Patronymic, dbUser.Patronymic)
	assert.Equal(t, user.Address, dbUser.Address)

	// Создание тестового HTTP-запроса для удаления пользователя
	urlDelete := fmt.Sprintf("/users/%d", dbUser.ID) // преобразуем dbUser.ID в строку
	reqDelete, _ := http.NewRequest("DELETE", urlDelete, nil)

	// Выполнение HTTP-запроса на удаление пользователя
	wDelete := httptest.NewRecorder()
	router.ServeHTTP(wDelete, reqDelete)

	// Проверка ответа на удаление пользователя
	assert.Equal(t, http.StatusOK, wDelete.Code)

	// Проверка, что пользователь был удален из базы данных
	var deletedUser models.User
	err = database.DB.Where("id = ?", dbUser.ID).First(&deletedUser).Error
	assert.Error(t, err, "expected error while fetching deleted user from database")
}
