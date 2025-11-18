package repository

import (
	"golang-app/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindByID(id uint) (*models.Order, error)
	Create(order *models.Order) error
	Update(order *models.Order) error
	CreateOrderItem(item *models.OrderItem) error
	FindOrdersByConsumer(consumerID uint, offset, limit int) ([]models.Order, int64, error)
	FindOrdersByAgent(agentID uint, offset, limit int) ([]models.Order, int64, error)
	FindOrderByAgent(orderID, agentID uint) (*models.Order, error)
	HasUserPurchasedProduct(userID, productID uint) (bool, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, id).Error
	return &order, err
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) CreateOrderItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) FindOrdersByConsumer(consumerID uint, offset, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	r.db.Model(&models.Order{}).Where("consumer_id = ?", consumerID).Count(&total)
	err := r.db.Preload("Items").Where("consumer_id = ?", consumerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) FindOrdersByAgent(agentID uint, offset, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	r.db.Model(&models.Order{}).Where("agent_id = ?", agentID).Count(&total)
	err := r.db.Preload("Items").Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) FindOrderByAgent(orderID, agentID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").Where("id = ? AND agent_id = ?", orderID, agentID).First(&order).Error
	return &order, err
}

func (r *orderRepository) HasUserPurchasedProduct(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Order{}).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.consumer_id = ? AND order_items.product_id = ? AND orders.status = ?", userID, productID, models.OrderStatusCompleted).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
