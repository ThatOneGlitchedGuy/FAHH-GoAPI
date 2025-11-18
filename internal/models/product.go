package models

import (
	"time"
)

type Product struct {
	ID          uint        `gorm:"primaryKey"`
	AgentID     uint        `gorm:"not null;index"`
	Title       string      `gorm:"not null;size:200"`
	Description string      `gorm:"type:text"`
	Price       float64     `gorm:"not null"`
	IsActive    bool        `gorm:"default:true"`
	CreatedAt   time.Time   `gorm:"not null"`
	Agent       User        `gorm:"foreignKey:AgentID"`
	OrderItems  []OrderItem `gorm:"foreignKey:ProductID"`
	Reviews     []Review    `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}