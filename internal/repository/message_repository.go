package repository

import (
	"golang-app/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *models.Message) error
	FindMessagesByOrderID(orderID uint, offset, limit int) ([]models.Message, int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) FindMessagesByOrderID(orderID uint, offset, limit int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64
	r.db.Model(&models.Message{}).Where("order_id = ?", orderID).Count(&total)
	err := r.db.Where("order_id = ?", orderID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, total, err
}
