package order

import "errors"

var ErrNotHandledStatus error = errors.New("not handled status")
var ErrNotFound error = errors.New("order not found")
