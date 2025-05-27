package main

import (
	"log"
	"os"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/subscription_service/controllers"
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

	subscriptionController := controllers.NewSubscriptionController()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	subscriptionRoutes := router.Group("/subscriptions")
	{
		subscriptionRoutes.POST("", subscriptionController.CreateSubscription)
		subscriptionRoutes.GET("", subscriptionController.GetSubscription)
		subscriptionRoutes.PATCH("", subscriptionController.UpdateSubscription)
		subscriptionRoutes.DELETE("/:id", subscriptionController.DeleteSubscription)
	}

	port := os.Getenv("SUBSCRIPTION_SERVICE_PORT")
	if port == "" {
		port = "8003"
	}

	log.Printf("Subscription service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
