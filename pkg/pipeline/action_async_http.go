package pipeline

import (
	"fmt"
	"os"
)

func init() {
	register("async-http", &actionAsyncHttp{})
}

type actionAsyncHttp struct{}

func (h *actionAsyncHttp) Execute(p *Pipeline, j *Job) (result *Result, err error) {
	f := availableActions["http"]

	fmt.Println(f)
	os.Exit(1)

	return nil, nil
}
