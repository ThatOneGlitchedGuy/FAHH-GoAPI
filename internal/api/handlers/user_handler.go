package handlers

import (
	"net/http"
	"strconv"

	"golang-app/internal/api/middleware"
	"golang-app/internal/schemas"
	"golang-app/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:id/profile", h.GetUserProfile)

	authed := router.Use(middleware.AuthMiddleware())
	{
		authed.GET("/me", h.GetCurrentUser)
		authed.PATCH("/me", h.UpdateCurrentUser)
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input schemas.UserBase
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := schemas.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(userID, &input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.service.GetUserProfile(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user profile"})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}