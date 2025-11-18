package service

import (
	"errors"
	"golang-app/internal/config"
	"golang-app/internal/models"
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
	"golang-app/internal/security"
	"math"
)

type MessageService interface {
	SendMessage(orderID, senderID uint, data *schemas.MessageCreate) (*schemas.MessageOut, error)
	ListMessages(orderID uint, userID uint, page, size int) (*schemas.Page, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	orderRepo   repository.OrderRepository
	config      *config.Settings
}

func NewMessageService(messageRepo repository.MessageRepository, orderRepo repository.OrderRepository, cfg *config.Settings) MessageService {
	return &messageService{messageRepo: messageRepo, orderRepo: orderRepo, config: cfg}
}

func (s *messageService) messageToMessageOut(message *models.Message) (*schemas.MessageOut, error) {
	decryptedContent, err := security.Decrypt(message.ContentEncrypted, s.config.FernetKey, 0)
	if err != nil {
		return nil, err
	}

	return &schemas.MessageOut{
		ID:        message.ID,
		OrderID:   message.OrderID,
		SenderID:  message.SenderID,
		Content:   string(decryptedContent),
		CreatedAt: message.CreatedAt,
	}, nil
}

func (s *messageService) SendMessage(orderID, senderID uint, data *schemas.MessageCreate) (*schemas.MessageOut, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.ConsumerID != senderID && order.AgentID != senderID {
		return nil, errors.New("not allowed to send messages for this order")
	}

	encryptedContent, err := security.Encrypt([]byte(data.Content), s.config.FernetKey)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		OrderID:          orderID,
		SenderID:         senderID,
		ContentEncrypted: encryptedContent,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return s.messageToMessageOut(message)
}

func (s *messageService) ListMessages(orderID uint, userID uint, page, size int) (*schemas.Page, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.ConsumerID != userID && order.AgentID != userID {
		return nil, errors.New("not allowed to view messages for this order")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = s.config.PageSizeDefault
	}
	if size > s.config.PageSizeMax {
		size = s.config.PageSizeMax
	}

	offset := (page - 1) * size
	messages, total, err := s.messageRepo.FindMessagesByOrderID(orderID, offset, size)
	if err != nil {
		return nil, err
	}

	messageOuts := make([]schemas.MessageOut, len(messages))
	for i, message := range messages {
		msgOut, err := s.messageToMessageOut(&message)
		if err != nil {
			return nil, err
		}
		messageOuts[i] = *msgOut
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := schemas.PageMeta{
		Page:  page,
		Size:  size,
		Total: total,
	}

	return &schemas.Page{
		Meta:  meta,
		Items: messageOuts,
	}, nil
}
