package handlers

import (
	"net/http"

	"golang-app/internal/api/middleware"
	"golang-app/internal/schemas"
	"golang-app/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authSvc service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{authService: authSvc, userService: userSvc}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.RegisterUser)
	router.POST("/login", h.LoginUser)
	router.GET("/me", middleware.AuthMiddleware(), h.GetCurrentUser)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input schemas.UserCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := schemas.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.RegisterUser(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, schemas.Token{AccessToken: token})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var input schemas.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := schemas.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, schemas.Token{AccessToken: token})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
