package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/plan_service/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanController struct {
	repo *repository.PlanRepository
}

func NewPlanController() *PlanController {
	return &PlanController{
		repo: repository.NewPlanRepository(),
	}
}

func (c *PlanController) CreatePlan(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Operation not permitted"})
		return
	}

	var req map[string]interface{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	var plan models.Plan

	if name, ok := req["name"].(string); ok {
		plan.Name = name
	}
	if price, ok := req["price"].(float64); ok {
		plan.Price = price
	}
	if features, ok := req["features"]; ok {
		b, _ := json.Marshal(features)
		var arr []string
		_ = json.Unmarshal(b, &arr)
		plan.Features = arr
	}
	if duration, ok := req["duration_days"].(float64); ok {
		plan.DurationDays = int(duration)
	}
	if meta, ok := req["meta"]; ok {
		b, _ := json.Marshal(meta)
		raw := json.RawMessage(b)
		plan.Meta = &raw
	}

	if err := models.ValidatePlan(&plan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	if err := c.repo.CreatePlan(&plan); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, plan)
}

func (c *PlanController) GetPlan(ctx *gin.Context) {
	var filter models.GetPlanFilter

	if idStr := ctx.Query("id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID format"})
			return
		}
		filter.ID = &id
	}
	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}
	if priceStr := ctx.Query("price"); priceStr != "" {
		var price float64
		_, err := fmt.Sscanf(priceStr, "%f", &price)
		if err == nil {
			filter.Price = &price
		}
	}
	if durationStr := ctx.Query("duration_days"); durationStr != "" {
		var duration int
		_, err := fmt.Sscanf(durationStr, "%d", &duration)
		if err == nil {
			filter.DurationDays = &duration
		}
	}
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	plans, err := c.repo.GetPlan(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plans: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, plans)
}

func (c *PlanController) UpdatePlan(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Operation not permitted"})
		return
	}

	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Plan ID is required"})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if features, ok := updateData["features"]; ok && features != nil {
		if arr, ok := features.([]interface{}); ok {
			strArr := make([]string, len(arr))
			for i, v := range arr {
				strArr[i] = fmt.Sprintf("%v", v)
			}
			updateData["features"] = strArr
		}
	}

	updatedPlan, err := c.repo.UpdatePlan(id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedPlan)
}

func (c *PlanController) DeletePlan(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Operation not permitted"})
		return
	}

	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Plan ID is required"})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID format"})
		return
	}

	if err := c.repo.DeletePlan(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
