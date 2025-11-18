package schemas

import (
	"time"
)

type MessageCreate struct {
	Content string `json:"content" validate:"required,min=1,max=4000"`
}

type MessageOut struct {
	ID        uint      `json:"id"`
	OrderID   uint      `json:"order_id"`
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}