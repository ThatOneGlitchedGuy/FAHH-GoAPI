package schemas

import (
	"time"
)

type ProductBase struct {
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"gte=0"`
	IsActive    bool    `json:"is_active" default:"true"`
}

type ProductCreate struct {
	ProductBase
}

type ProductUpdate struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

type ProductOut struct {
	ID          uint      `json:"id"`
	AgentID     uint      `json:"agent_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}