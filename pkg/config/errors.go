package config

import "errors"

var (
	ErrUnknownAPIVersion error = errors.New("unknown api version")
	ErrUnknownKind       error = errors.New("unknown kind")
)
