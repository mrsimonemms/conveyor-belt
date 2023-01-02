package server

import "errors"

var (
	ErrNoTriggers     error = errors.New("no triggers configured")
	ErrUnknownTrigger error = errors.New("unknown trigger")
)
