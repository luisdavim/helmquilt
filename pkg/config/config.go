package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const DefaultConfigFile = "./helmquilt.yaml"

var ErrMissingConfig = errors.New("missing config file")

type Options struct {
	Force      bool
	Repack     bool
	ConfigFile string
	WorkDir    string
}

func (opts *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force run (ignore lock file)")
	cmd.Flags().BoolVarP(&opts.Repack, "repack", "r", false, "Repack the chart as a tarball")
	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", DefaultConfigFile, "path to the config file")
}

type Config struct {
	Charts []Chart `json:"charts"`
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
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(opts.ConfigFile, out, 0o644)
}
