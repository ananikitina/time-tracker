package migrations

import (
	"time-tracker/database"
	"time-tracker/models"
)

func Migrate() {
	database.DB.AutoMigrate(&models.User{}, &models.Task{})
}
