package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

type Options struct {
	Force      bool
	Repack     bool
	ConfigFile string
	WorkDir    string
}

type Config struct {
	Charts []Chart `json:"charts"`
}

func Read(configFile string) (Config, error) {
	var charts Config

	if _, err := os.Stat(configFile); err != nil {
		return charts, fmt.Errorf("missing config file: %w", err)
	}
	data, _ := os.ReadFile(configFile)
	if err := yaml.Unmarshal(data, &charts); err != nil {
		return charts, err
	}

	return charts, nil
}

func Save(cfg Config, opts Options) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(opts.ConfigFile, out, 0o644)
}
