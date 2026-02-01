package courier

import "errors"

var ErrConflict error = errors.New("conflict")
var ErrNotFound error = errors.New("not found")
var ErrNothingToUpdate = errors.New("nothing to update")
