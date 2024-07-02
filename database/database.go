package database

import (
	"fmt"
	"log"
	"os"

	"time-tracker/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Объявляем глобальную переменную DB типа *gorm.DB для использования в других частях программы.
var DB *gorm.DB

// Функция Connect выполняет подключение к базе данных и миграцию схемы.
func Connect() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Формируем строку подключения к базе данных, используя переменные окружения.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	// Открываем соединение с базой данных PostgreSQL с использованием GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	// Присваиваем открытую базу данных глобальной переменной DB.
	DB = db

	// Выполняем автоматическую миграцию схемы базы данных для моделей User и Task.
	db.AutoMigrate(&models.User{}, &models.Task{})
}
