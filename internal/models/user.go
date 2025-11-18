package models

import (
	"time"
)

type User struct {
	ID                uint      `gorm:"primaryKey;index"`
	Email             string    `gorm:"uniqueIndex;not null;size:255"`
	HashedPassword    string    `gorm:"not null;size:255"`
	Role              UserRole  `gorm:"type:varchar(50);default:'CONSUMER';not null"`
	FullNameEncrypted []byte    `gorm:"type:BLOB"`
	AddressEncrypted  []byte    `gorm:"type:BLOB"`
	CreatedAt         time.Time `gorm:"not null"`
	IsActive          bool      `gorm:"default:true"`
	Products          []Product `gorm:"foreignKey:AgentID"`
	OrdersAsConsumer  []Order   `gorm:"foreignKey:ConsumerID"`
	OrdersAsAgent     []Order   `gorm:"foreignKey:AgentID"`
	Reviews           []Review  `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}