package models

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"default:1"`
	UnitPrice float64 `gorm:"not null"`
	Order     Order   `gorm:"foreignKey:OrderID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}