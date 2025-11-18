package service

import (
	"errors"
	"golang-app/internal/config"
	"golang-app/internal/models"
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
	"math"
)

type OrderService interface {
	CreateOrder(data *schemas.OrderCreate, consumerID uint) (*schemas.OrderOut, error)
	ListMyOrders(userID uint, userRole models.UserRole, page, size int) (*schemas.Page, error)
	GetOrderByID(orderID uint, userID uint) (*schemas.OrderOut, error)
	UpdateOrderStatus(orderID uint, data *schemas.OrderUpdateStatus, agentID uint) (*schemas.OrderOut, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	config      *config.Settings
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, cfg *config.Settings) OrderService {
	return &orderService{orderRepo: orderRepo, productRepo: productRepo, config: cfg}
}

func (s *orderService) orderToOrderOut(order *models.Order) *schemas.OrderOut {
	orderItemsOut := make([]schemas.OrderItemOut, len(order.Items))
	for i, item := range order.Items {
		orderItemsOut[i] = schemas.OrderItemOut{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return &schemas.OrderOut{
		ID:          order.ID,
		ConsumerID:  order.ConsumerID,
		AgentID:     order.AgentID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		Items:       orderItemsOut,
	}
}

func (s *orderService) CreateOrder(data *schemas.OrderCreate, consumerID uint) (*schemas.OrderOut, error) {
	if len(data.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	productIDs := make([]uint, len(data.Items))
	for i, item := range data.Items {
		productIDs[i] = item.ProductID
	}

	products, err := s.productRepo.FindProductsByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productsMap := make(map[uint]models.Product)
	for _, p := range products {
		productsMap[p.ID] = p
	}

	if len(productsMap) != len(productIDs) {
		return nil, errors.New("one or more products not found or inactive")
	}

	agentIDs := make(map[uint]struct{})
	for _, item := range data.Items {
		product, ok := productsMap[item.ProductID]
		if !ok {
			return nil, errors.New("product not found")
		}
		agentIDs[product.AgentID] = struct{}{}
	}

	if len(agentIDs) != 1 {
		return nil, errors.New("all products in an order must belong to the same agent")
	}
	var agentID uint
	for id := range agentIDs {
		agentID = id
	}

	order := &models.Order{
		ConsumerID:  consumerID,
		AgentID:     agentID,
		Status:      models.OrderStatusNew,
		TotalAmount: 0.0,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, itemInput := range data.Items {
		product := productsMap[itemInput.ProductID]
		unitPrice := product.Price
		total += unitPrice * float64(itemInput.Quantity)

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: itemInput.ProductID,
			Quantity:  itemInput.Quantity,
			UnitPrice: unitPrice,
		}
		// Assuming orderRepo can handle order items or a separate OrderItemRepo
		// For simplicity, directly creating order items here.
		if err := s.orderRepo.CreateOrderItem(&orderItem); err != nil {
			return nil, err
		}
	}

	order.TotalAmount = math.Round(total*100) / 100
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	// Reload order with items
	order, err = s.orderRepo.FindByID(order.ID)
	if err != nil {
		return nil, err
	}

	return s.orderToOrderOut(order), nil
}

func (s *orderService) ListMyOrders(userID uint, userRole models.UserRole, page, size int) (*schemas.Page, error) {
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
	var orders []models.Order
	var total int64
	var err error

	if userRole == models.UserRoleConsumer {
		orders, total, err = s.orderRepo.FindOrdersByConsumer(userID, offset, size)
	} else {
		orders, total, err = s.orderRepo.FindOrdersByAgent(userID, offset, size)
	}

	if err != nil {
		return nil, err
	}

	ordersOut := make([]schemas.OrderOut, len(orders))
	for i, order := range orders {
		ordersOut[i] = *s.orderToOrderOut(&order)
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
		Items: ordersOut,
	}, nil
}

func (s *orderService) GetOrderByID(orderID uint, userID uint) (*schemas.OrderOut, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.ConsumerID != userID && order.AgentID != userID {
		return nil, errors.New("not allowed to access this order")
	}

	return s.orderToOrderOut(order), nil
}

func (s *orderService) UpdateOrderStatus(orderID uint, data *schemas.OrderUpdateStatus, agentID uint) (*schemas.OrderOut, error) {
	order, err := s.orderRepo.FindOrderByAgent(orderID, agentID)
	if err != nil {
		return nil, err
	}

	allowedTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusNew:      {models.OrderStatusAccepted, models.OrderStatusRejected, models.OrderStatusCanceled},
		models.OrderStatusAccepted: {models.OrderStatusShipped, models.OrderStatusCanceled},
		models.OrderStatusShipped:  {models.OrderStatusCompleted},
	}

	allowedNextStatuses, ok := allowedTransitions[order.Status]
	if !ok {
		return nil, errors.New("invalid current order status for transition")
	}

	isValidTransition := false
	for _, s := range allowedNextStatuses {
		if s == data.Status {
			isValidTransition = true
			break
		}
	}

	if !isValidTransition {
		return nil, errors.New("invalid status transition")
	}

	order.Status = data.Status
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderToOrderOut(order), nil
}
