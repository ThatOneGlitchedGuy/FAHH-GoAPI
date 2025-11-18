package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	UserID    uint   `gorm:"not null;index"`
	Rating    int    `gorm:"not null"`
	Comment   string `gorm:"type:text"`
	Product   Product
	User      User
}
