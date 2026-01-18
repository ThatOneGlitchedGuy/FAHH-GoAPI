package handlers

import (
	"net/http"
	"strconv"

	"golang-app/internal/api/middleware"
	"golang-app/internal/models"
	"golang-app/internal/schemas"
	"golang-app/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	authed := router.Use(middleware.AuthMiddleware())
	{
		authed.POST("/", middleware.OnlyConsumer(), h.CreateOrder)
		authed.GET("/", h.ListMyOrders)
		authed.GET("/:order_id", h.GetOrder)
		authed.PATCH("/:order_id/status", middleware.OnlyAgent(), h.UpdateOrderStatus)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	consumerID := c.MustGet("userID").(uint)

	var input schemas.OrderCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := schemas.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(&input, consumerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(models.UserRole)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	ordersPage, err := h.service.ListMyOrders(userID, userRole, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, ordersPage)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.service.GetOrderByID(uint(orderID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	agentID := c.MustGet("userID").(uint)

	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var input schemas.OrderUpdateStatus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := schemas.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.UpdateOrderStatus(uint(orderID), &input, agentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or not owned by agent"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
