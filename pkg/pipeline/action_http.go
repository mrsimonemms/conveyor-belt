package pipeline

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func init() {
	register("http", &actionHttp{})
}

type actionHttp struct{}

func (h *actionHttp) Execute(p *Pipeline, j *Job) (result *Result, err error) {
	var timeout time.Duration

	if j.Timeout == nil {
		// No timeout specified - set to 30 seconds
		timeout = 30 * time.Second
	} else {
		timeout = *j.Timeout
	}

	// Create an HTTP client with the appropriate timeout
	client := http.Client{
		Timeout: timeout,
	}

	// Convert data to JSON
	// @todo(sje): support multiple body data types
	data, err := parseData(p, &j.Action.HTTP.Data)
	if err != nil {
		return
	}

	// Convert the data to io.Reader
	requestBody := bytes.NewBuffer(data)

	req, err := http.NewRequest(strings.ToUpper(j.Action.HTTP.Method), j.Action.HTTP.URL, requestBody)
	if err != nil {
		return
	}

	// Set the content type
	// @todo(sje): support multiple content types
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		err = res.Body.Close()
	}()
	if err != nil {
		return
	}

	result = &Result{
		Status:  res.StatusCode,
		Headers: res.Header,
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	// @todo(sje): support multiple content types
	err = json.Unmarshal(body, &result.Body)
	if err != nil {
		return
	}

	return
}
