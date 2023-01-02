package pipeline

import (
	"fmt"
	"net/http"

	"github.com/mrsimonemms/conveyor-belt/pkg/config"
)

type Job struct {
	config.PipelineJob
}

func (j *Job) getAction() (Action, error) {
	key := j.Action.GetActionKey()

	if key == nil {
		return nil, ErrActionUnknown
	}

	action, ok := availableActions[*key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrActionNotRegistered, *key)
	}

	return action, nil
}

func (j *Job) Exec(p *Pipeline) (*Result, error) {
	// Select the desired action
	action, err := j.getAction()
	if err != nil {
		return nil, ErrActionUnknown
	}

	return action.Execute(p, j)
}

type Result struct {
	Status  int
	Headers http.Header
	Body    map[string]any
}
