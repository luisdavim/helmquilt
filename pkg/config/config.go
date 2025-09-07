package config

import (
	"errors"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

const DefaultConfigFile = "./helmquilt.yaml"

var ErrMissingConfig = errors.New("missing config file")

type Config struct {
	Charts []Chart `json:"charts"`
}

func (c *Config) SetDefaults() error {
	for _, chart := range c.Charts {
		if err := chart.SetDefaults(); err != nil {
			return err
		}
	}

	return nil
}

func Read(configFile string) (Config, error) {
	var charts Config

	if _, err := os.Stat(configFile); err != nil {
		return charts, fmt.Errorf("%w: %w", ErrMissingConfig, err)
	}

	data, _ := os.ReadFile(configFile)
	if err := yaml.Unmarshal(data, &charts); err != nil {
		return charts, err
	}

	return charts, nil
}

func Save(cfg Config, opts Options) error {
	if opts.DryRun {
		return nil
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(opts.ConfigFile, out, 0o644)
}
