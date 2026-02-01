package courier

type CourierResponse struct {
	ID            *int64  `json:"id"`
	Name          *string `json:"name"`
	Phone         *string `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}

type RequestCreate struct {
	Name          *string `json:"name"`
	Phone         *string `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}

type RequestUpdate CourierResponse
