package delivery

import "errors"

var ErrNotFound error = errors.New("not found")
var ErrConflict error = errors.New("conflict")
var ErrCannotCalculateDeliveryTime error = errors.New("cannot calculate delivery time")
