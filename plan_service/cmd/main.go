package main

import (
	"log"
	"os"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/auth"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/plan_service/controllers"
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
	planController := controllers.NewPlanController()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	admin := router.Group("")
	admin.Use(auth.AuthMiddleware())
	{
		admin.GET("/plans", planController.GetPlan)
		admin.POST("/plans", planController.CreatePlan)
		admin.PATCH("/plans", planController.UpdatePlan)
		admin.DELETE("/plans/:id", planController.DeletePlan)
	}

	port := os.Getenv("PLAN_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Plan service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
