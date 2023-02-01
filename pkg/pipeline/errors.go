package pipeline

import "errors"

var (
	ErrActionRegistered     error = errors.New("action already registered")
	ErrActionNotRegistered  error = errors.New("required action not registered")
	ErrActionUnknown        error = errors.New("action unknown")
	ErrNoJobs               error = errors.New("no jobs configured")
	ErrUnknownStage         error = errors.New("unknown stage")
	ErrUnknownAsyncCallback error = errors.New("unknown async callback")
)
