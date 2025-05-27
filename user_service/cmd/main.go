package main

import (
	"log"
	"os"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/auth"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/user_service/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	database.Connect()
	defer database.Close()

	router := gin.Default()

	userController := controllers.NewUserController()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	router.POST("/login", userController.Login)
	router.POST("/users", userController.CreateUser) // would require 2FA to be protected

	protected := router.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/users", userController.GetUser)
		protected.PATCH("/users", userController.UpdateUser)
		protected.DELETE("/users/:id", userController.DeleteUser)
	}

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}