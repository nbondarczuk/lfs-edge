package files

import (
	"errors"
)

var (
	ErrInvalidArg  = errors.New("the passed argument format is invalid")
	ErrNotFound    = errors.New("the requested file was not found")
	ErrNotAllowed  = errors.New("the requested access is not allowed")	
)
