package rest

import (
	"errors"
)

var (
	ErrUrlDecode             = errors.New("error decoding URL path")
	ErrPayloadRead           = errors.New("payload read error")
	ErrAuthz                 = errors.New("authorization error")
	ErrPathVariableMissing   = errors.New("a required path variable was not specified")
	ErrQueryParameterMissing = errors.New("a required quesry paramter was not specified")
)
