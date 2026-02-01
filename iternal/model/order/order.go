package order

import "time"

type Order struct {
	ID        string
	Status    Status
	CreatedAt time.Time
}

type Status string

const (
	OrderStatusPending    Status = "pending"
	OrderStatusConfirmed  Status = "confirmed"
	OrderStatusCooking    Status = "cooking"
	OrderStatusDelivering Status = "delivering"
	OrderStatusDelivered  Status = "delivered"
	OrderStatusCanceled   Status = "canceled"
	OrderStatusDeleted    Status = "deleted"
	OrderStatusCreated    Status = "created"
	OrderStatusUpdated    Status = "updated"
	OrderStatusCompleted  Status = "completed"
)

var ValidStatuses = map[string]Status{
	"pending":    OrderStatusPending,
	"confirmed":  OrderStatusConfirmed,
	"cooking":    OrderStatusCooking,
	"delivering": OrderStatusDelivering,
	"delivered":  OrderStatusDelivered,
	"canceled":   OrderStatusCanceled,
	"deleted":    OrderStatusDeleted,
	"created":    OrderStatusCreated,
	"updated":    OrderStatusUpdated,
	"completed":  OrderStatusCompleted,
}
