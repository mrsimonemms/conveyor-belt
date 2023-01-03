package config

import (
	"time"

	"k8s.io/utils/pointer"
)

type Config struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   Metadata       `json:"metadata"`
	Spec       PipelineConfig `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type PipelineConfig struct {
	Error    *PipelineJob      `json:"error,omitempty"`
	Jobs     []PipelineJob     `json:"jobs"`
	Port     int               `json:"port"`
	Stages   []string          `json:"stages"`
	Triggers []PipelineTrigger `json:"triggers"`
}

type PipelineAction struct {
	HTTP *PipelineActionHTTP `json:"http,omitempty"`
}

func (a PipelineAction) GetActionKey() *string {
	if a.HTTP != nil {
		return pointer.String("http")
	}

	return nil
}

// @todo(sje): might be able to use something else
type PipelineActionHTTP struct {
	Method string         `json:"method"`
	URL    string         `json:"url"`
	Data   map[string]any `json:"data"`
}

type PipelineJob struct {
	Name    string         `json:"name"`
	Stage   string         `json:"stage"`
	Action  PipelineAction `json:"action"`
	Timeout *time.Duration `json:"timeout,omitempty"`
}

type PipelineTrigger struct {
	Type PipelineTriggerType `json:"type"`
}

type PipelineTriggerType string

const (
	PipelineTriggerTypeWebHook PipelineTriggerType = "webhook"
)

// @todo(sje): validate the config
func (c *Config) IsValid() error {
	return nil
}
