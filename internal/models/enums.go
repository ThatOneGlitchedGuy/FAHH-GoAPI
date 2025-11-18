package models

type UserRole string

const (
	UserRoleAgent    UserRole = "AGENT"
	UserRoleConsumer UserRole = "CONSUMER"
	UserRoleAdmin    UserRole = "ADMIN"
)

type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "NEW"
	OrderStatusAccepted  OrderStatus = "ACCEPTED"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCanceled  OrderStatus = "CANCELED"
)