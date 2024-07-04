package routes

import (
	"time-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Define a route group for users
	userRoutes := r.Group("/users")
	{
		// Routes related to general user management
		userRoutes.GET("", handlers.GetUsers)          // Fetch users with filtering and pagination.
		userRoutes.DELETE("/:id", handlers.DeleteUser) // Delete a user by ID.
		userRoutes.PUT("/:id", handlers.UpdateUser)    // Update user data by ID.
		userRoutes.POST("", handlers.AddUser)          // Add a new user.

		// Routes related to tasks and efforts for a specific user
		userRoutes.GET("/:id/tasks", handlers.GetUserTasks) // Get tasks for a specific user.
		//userRoutes.GET("/efforts/:id", handlers.GetUserEffort)  // Get efforts for a specific user.
		userRoutes.POST("/:userID/tasks/start", handlers.StartTask) // Start tracking time for a task for a user.
		//userRoutes.PUT("/:userID/tasks/finish", handlers.FinishTask) // Finish tracking time for a task for a user.
	}
}
