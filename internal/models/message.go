package models

import (
	"time"
)

type Message struct {
	ID               uint      `gorm:"primaryKey"`
	OrderID          uint      `gorm:"not null;index"`
	SenderID         uint      `gorm:"not null"`
	ContentEncrypted []byte    `gorm:"type:BLOB;not null"`
	CreatedAt        time.Time `gorm:"not null"`
	Order            Order     `gorm:"foreignKey:OrderID"`
}

func (Message) TableName() string {
	return "messages"
}