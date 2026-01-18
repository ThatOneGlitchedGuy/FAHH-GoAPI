package handlers

import (
	"net/http"

	"golang-app/internal/api/middleware"
	"golang-app/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service service.StatsService
}

func NewAdminHandler(svc service.StatsService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.OnlyAdmin())
	{
		adminGroup.GET("/stats", h.GetStats)
	}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetAdminStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
