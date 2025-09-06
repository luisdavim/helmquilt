package config

import "github.com/spf13/cobra"

type Options struct {
	Force      bool
	Repack     bool
	ConfigFile string
	WorkDir    string
	DryRun     bool
}

func (opts *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force run (ignore lock file)")
	cmd.Flags().BoolVarP(&opts.Repack, "repack", "r", false, "Repack the chart as a tarball")
	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", DefaultConfigFile, "path to the config file")
	cmd.Flags().StringVarP(&opts.WorkDir, "workdir", "W", "", "Override workdir, instead of config file location")
}

type DiffOptions struct {
	Options
	Write bool
}

func (ops *DiffOptions) AddFlags(cmd *cobra.Command) {
	ops.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&ops.Write, "write", "w", false, "Write patch files")
}
