package controllers

import (
	"net/http"
	"time"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/utils"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/subscription_service/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionController struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{
		repo: repository.NewSubscriptionRepository(),
	}
}

func (c *SubscriptionController) CreateSubscription(ctx *gin.Context) {
	var subscription models.Subscription
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateUserExists(subscription.UserID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidatePlanExists(subscription.PlanID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateSubscription(&subscription); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.CreateSubscription(&subscription); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, subscription)
}

func (c *SubscriptionController) GetSubscription(ctx *gin.Context) {
	var filter models.GetSubscriptionFilter

	if idStr := ctx.Query("id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID format"})
			return
		}
		filter.ID = &id
	}

	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		filter.UserID = &userID
	}

	if planIDStr := ctx.Query("plan_id"); planIDStr != "" {
		planID, err := uuid.Parse(planIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID format"})
			return
		}
		filter.PlanID = &planID
	}

	subscriptions, err := c.repo.GetSubscription(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, sub := range subscriptions {
		if sub.Status != "EXPIRED" && time.Now().After(sub.EndDate) {
			updates := map[string]interface{}{
				"status": "EXPIRED",
			}
			_, err := c.repo.UpdateSubscription(sub.ID, updates)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription status"})
				return
			}
			sub.Status = "EXPIRED"
		}
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

func (c *SubscriptionController) UpdateSubscription(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID format"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if planIDInterface, ok := updates["plan_id"]; ok {
		if planIDStr, ok := planIDInterface.(string); ok {
			if err := utils.ValidatePlanExists(uuid.MustParse(planIDStr)); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if userIDInterface, ok := updates["user_id"]; ok {
		if userIDStr, ok := userIDInterface.(string); ok {
			if err := utils.ValidateUserExists(uuid.MustParse(userIDStr)); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	updatedSubscription, err := c.repo.UpdateSubscription(id, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedSubscription)
}

func (c *SubscriptionController) DeleteSubscription(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID format"})
		return
	}

	if err := c.repo.DeleteSubscription(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
