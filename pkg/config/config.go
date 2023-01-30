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
	Async    *PipelineAsync    `yaml:"async,omitempty"`
	Error    *PipelineJob      `yaml:"error,omitempty"`
	Jobs     []PipelineJob     `yaml:"jobs"`
	Port     int               `yaml:"port"`
	Stages   []string          `yaml:"stages"`
	Triggers []PipelineTrigger `yaml:"triggers"`
}

type PipelineAsync struct {
	Domain string `yaml:"domain"`
}

type PipelineAction struct {
	HTTP *PipelineActionHTTP `yaml:"http,omitempty"`
}

func (a PipelineAction) GetActionKey() *string {
	if a.HTTP != nil {
		return pointer.String("http")
	}

	return nil
}

// @todo(sje): might be able to use something else
type PipelineActionHTTP struct {
	Method string         `yaml:"method"`
	URL    string         `yaml:"url"`
	Data   map[string]any `yaml:"data"`
}

type PipelineJob struct {
	Async   *PipelineJobAsync `yaml:"async,omitempty"`
	Name    string            `yaml:"name"`
	Stage   string            `yaml:"stage"`
	Action  PipelineAction    `yaml:"action"`
	Timeout *time.Duration    `yaml:"timeout,omitempty"`
}

type PipelineJobAsync struct {
	Enabled           bool   `yaml:"enabled"`
	CallbackURLHeader string `yaml:"callbackUrlHeader"`
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
