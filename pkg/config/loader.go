package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

func Load(filepath string) (*Config, error) {
	cfgFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	if err := yaml.Unmarshal(cfgFile, &cfg); err != nil {
		return nil, err
	}

	// @todo(sje): support multiple api versions and kinds
	if cfg.APIVersion != "conveyor-belt.simonemms.com/v1alpha1" {
		return nil, ErrUnknownAPIVersion
	}
	if cfg.Kind != "Pipeline" {
		return nil, ErrUnknownKind
	}

	if err := cfg.IsValid(); err != nil {
		return nil, err
	}

	return cfg, nil
}
