package routes

import (
	"github.com/ananikitina/time-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", handlers.GetUsers)
		userRoutes.DELETE("/:id", handlers.DeleteUser)
		userRoutes.PUT("/:id", handlers.UpdateUser)
		userRoutes.POST("", handlers.AddUser)
		userRoutes.GET("/:id", handlers.GetUserTasks)
	}
	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.GET("/:userID/sort", handlers.SortTasks)
		taskRoutes.POST("/:userID/start", handlers.StartTask)
		taskRoutes.PUT("/:userID/finish", handlers.FinishTask)
	}
}
