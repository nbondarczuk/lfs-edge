package file

// Common definition of the file status to be used by proxy and sync.
// As the values are used in db they can't be internal.
const (
	StatusOk       = 0
	StatusNotFound = 1
	StatusPending  = 2
	StatusUnknown  = 3
	StatusError    = 4
)
