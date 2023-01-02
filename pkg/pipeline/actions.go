package pipeline

import (
	"fmt"
)

type Action interface {
	Execute(*Pipeline, *Job) (*Result, error)
}

var availableActions map[string]Action

func register(name string, action Action) {
	if availableActions == nil {
		availableActions = make(map[string]Action, 0)
	}

	if _, ok := availableActions[name]; ok {
		panic(fmt.Errorf("%s: %s", ErrActionRegistered, name))
	}

	availableActions[name] = action
}
