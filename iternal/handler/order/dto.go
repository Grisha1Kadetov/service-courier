package order

import "time"

type Order struct {
	ID        string    `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
