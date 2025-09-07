package config

import "github.com/spf13/cobra"

type Options struct {
	ConfigFile string
	WorkDir    string
	DryRun     bool
}

func (opts *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", DefaultConfigFile, "path to the config file")
	cmd.Flags().StringVarP(&opts.WorkDir, "workdir", "W", "", "Override workdir, instead of config file location")
}

type ApplyOptions struct {
	Options
	Force  bool
	Repack bool
}

func (opts *ApplyOptions) AddFlags(cmd *cobra.Command) {
	opts.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force run (ignore lock file)")
	cmd.Flags().BoolVarP(&opts.Repack, "repack", "r", false, "Repack the chart as a tarball")
}

type DiffOptions struct {
	Options
	Write bool
}

func (opts *DiffOptions) AddFlags(cmd *cobra.Command) {
	opts.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&opts.Write, "write", "w", false, "Write patch files")
}
