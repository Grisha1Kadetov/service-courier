package delivery

import "time"

type AssignResponse struct {
	CourierID     *int64     `json:"courier_id"`
	OrderID       *string    `json:"order_id"`
	TransportType *string    `json:"transport_type"`
	Deadline      *time.Time `json:"delivery_deadline"`
}

type UnassignResponse struct {
	OrderID   *string `json:"order_id"`
	Status    *string `json:"status"`
	CourierID *int64  `json:"courier_id"`
}

type AssignRequest struct {
	OrderID *string `json:"order_id"`
}

type UnassignRequest AssignRequest
