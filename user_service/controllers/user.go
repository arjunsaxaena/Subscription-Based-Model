package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/auth"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/user_service/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo *repository.UserRepository
}

func NewUserController() *UserController {
	return &UserController{
		repo: repository.NewUserRepository(),
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req map[string]interface{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	var user models.User

	if email, ok := req["email"].(string); ok {
		user.Email = email
	}
	if password, ok := req["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password: " + err.Error()})
			return
		}
		user.PasswordHash = string(hashedPassword)
	}
	if name, ok := req["name"].(string); ok {
		user.Name = &name
	}
	if meta, ok := req["meta"]; ok {
		b, _ := json.Marshal(meta)
		raw := json.RawMessage(b)
		user.Meta = &raw
	}

	if err := models.ValidateUser(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	if err := c.repo.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		userID := ctx.MustGet("userID").(uuid.UUID)
		filter := models.GetUserFilter{
			ID: &userID,
		}
		users, err := c.repo.GetUser(filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, users)
		return
	}

	var filter models.GetUserFilter

	if idStr := ctx.Query("id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		filter.ID = &id
	}

	if email := ctx.Query("email"); email != "" {
		filter.Email = &email
	}

	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}

	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	users, err := c.repo.GetUser(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		authenticatedUserID := ctx.MustGet("userID").(uuid.UUID)
		idStr := ctx.Query("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		if id != authenticatedUserID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
			return
		}
	}

	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if newEmail, ok := updateData["email"].(string); ok {
		if newEmail == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email cannot be empty"})
			return
		}
		filter := models.GetUserFilter{Email: &newEmail}
		users, err := c.repo.GetUser(filter)
		if err == nil && len(users) > 0 {
			for _, u := range users {
				if u.ID != id {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
					return
				}
			}
		}
	}

	if password, ok := updateData["password"]; ok {
		if password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password: " + err.Error()})
			return
		}
		updateData["password_hash"] = string(hashedPassword)
		delete(updateData, "password")
	}

	updatedUser, err := c.repo.UpdateUser(id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	isInternalServer := ctx.GetBool("isInternalServer")
	if !isInternalServer {
		authenticatedUserID := ctx.MustGet("userID").(uuid.UUID)
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		if id != authenticatedUserID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own profile"})
			return
		}
	}

	idStr := ctx.Param("id")

	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := c.repo.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *UserController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := models.GetUserFilter{
		Email: &req.Email,
	}
	users, err := c.repo.GetUser(filter)
	if err != nil || len(users) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}