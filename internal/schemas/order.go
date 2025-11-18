package schemas

import (
	"time"

	"golang-app/internal/models"
)

type OrderItemIn struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gte=1"`
}

type OrderCreate struct {
	Items []OrderItemIn `json:"items" validate:"required,min=1,dive"`
}

type OrderItemOut struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderOut struct {
	ID          uint               `json:"id"`
	ConsumerID  uint               `json:"consumer_id"`
	AgentID     uint               `json:"agent_id"`
	Status      models.OrderStatus `json:"status"`
	TotalAmount float64            `json:"total_amount"`
	CreatedAt   time.Time          `json:"created_at"`
	Items       []OrderItemOut     `json:"items"`
}

type OrderUpdateStatus struct {
	Status models.OrderStatus `json:"status" validate:"required"`
}