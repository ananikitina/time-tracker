package main

import (
	"log"

	"github.com/ananikitina/time-tracker/database"
	"github.com/ananikitina/time-tracker/routes"

	_ "github.com/ananikitina/time-tracker/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Time Tracker API
// @version 1.0
// @description This is a sample time tracker server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.Connect()
	database.Migrate()

	// Gin initialization
	r := gin.Default()

	// Routes registration
	routes.SetupRouter(r)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Starting the server
	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
