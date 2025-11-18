package models

import (
	"time"
)

type Order struct {
	ID          uint        `gorm:"primaryKey"`
	ConsumerID  uint        `gorm:"not null;index"`
	AgentID     uint        `gorm:"not null;index"`
	Status      OrderStatus `gorm:"type:varchar(50);default:'NEW';not null"`
	TotalAmount float64     `gorm:"default:0.0"`
	CreatedAt   time.Time   `gorm:"not null"`
	Consumer    User        `gorm:"foreignKey:ConsumerID"`
	Agent       User        `gorm:"foreignKey:AgentID"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	Messages    []Message   `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}