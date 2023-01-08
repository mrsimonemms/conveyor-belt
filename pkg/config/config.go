package config

import (
	"time"

	"k8s.io/utils/pointer"
)

type Config struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       PipelineConfig `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

type PipelineConfig struct {
	Error    *PipelineJob      `yaml:"error,omitempty"`
	Jobs     []PipelineJob     `yaml:"jobs"`
	Port     int               `yaml:"port"`
	Stages   []string          `yaml:"stages"`
	Triggers []PipelineTrigger `yaml:"triggers"`
}

type PipelineAction struct {
	HTTP      *PipelineActionHTTP      `yaml:"http,omitempty"`
	AsyncHTTP *PipelineActionAsyncHTTP `yaml:"asyncHttp,omitempty"`
}

func (a PipelineAction) GetActionKey() *string {
	if a.HTTP != nil {
		return pointer.String("http")
	}
	if a.AsyncHTTP != nil {
		return pointer.String("async-http")
	}

	return nil
}

// @todo(sje): might be able to use something else
type PipelineActionHTTP struct {
	Method string         `yaml:"method"`
	URL    string         `yaml:"url"`
	Data   map[string]any `yaml:"data"`
}

type PipelineActionAsyncHTTP struct {
	PipelineActionHTTP `yaml:",inline"`
}

type PipelineJob struct {
	Name    string         `yaml:"name"`
	Stage   string         `yaml:"stage"`
	Action  PipelineAction `yaml:"action"`
	Timeout *time.Duration `yaml:"timeout,omitempty"`
}

type PipelineTrigger struct {
	Type PipelineTriggerType `yaml:"type"`
}

type PipelineTriggerType string

const (
	PipelineTriggerTypeWebHook PipelineTriggerType = "webhook"
)

// @todo(sje): validate the config
func (c *Config) IsValid() error {
	return nil
}
