package courier

import "time"

type Courier struct {
	ID            *int64
	Name          *string
	Phone         *string
	Status        *Status
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	TransportType *TransportType
}

type Status string

const (
	StatusAvalible Status = "available"
	StatusBusy     Status = "busy"
	StatusPaused   Status = "paused"
)

type TransportType string

const (
	TransportScooter TransportType = "scooter"
	TransportCar     TransportType = "car"
	TransportOnFoot  TransportType = "on_foot"
)
