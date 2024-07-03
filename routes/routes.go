package routes

import (
	"time-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.GET("/users", handlers.GetUsers)                       // Fetch users with filtering and pagination.
	r.GET("/users/:id/task-efforts", handlers.GetTaskEffort) // Retrieve task efforts for a user within a specified period, sorted by effort (descending).
	r.POST("/users/:id/start-task", handlers.StartTask)      // Start tracking time for a task for a user.
	r.POST("/users/:id/stop-task", handlers.StopTask)        // Stop tracking time for a task for a user.
	r.DELETE("/users/:id", handlers.DeleteUser)              // Delete a user.
	r.PUT("/users/:id", handlers.UpdateUser)                 // Update user data.
	r.POST("/user", handlers.AddUser)                        // Add a new user.

}
