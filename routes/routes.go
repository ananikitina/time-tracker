package routes

import (
	"time-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()            // Создание нового маршрутизатора с настройками по умолчанию
	userGroup := r.Group("/user") // Создание группы маршрутов для URL, начинающихся с /users
	{
		userGroup.GET("/", handlers.GetUsers)
		userGroup.POST("/", handlers.AddUser)
	}
	return r
}
