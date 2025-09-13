package config

import "github.com/spf13/cobra"

// Options represents the genral/common CLI options
type Options struct {
	ConfigFile string
	WorkDir    string
	DryRun     bool
	Quiet      bool
}

// AddFlags adds the flags for this set of options to the given command
func (opts *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", DefaultConfigFile, "path to the config file")
	cmd.Flags().StringVarP(&opts.WorkDir, "workdir", "W", "", "Override workdir, instead of config file location")
}

// ApplyOptions represents the CLI options for the apply command
type ApplyOptions struct {
	Options
	Force  bool
	Repack bool
}

// AddFlags adds the flags for this set of options to the given command
func (opts *ApplyOptions) AddFlags(cmd *cobra.Command) {
	opts.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force run (ignore lock file)")
	cmd.Flags().BoolVarP(&opts.Repack, "repack", "r", false, "Repack the chart as a tarball")
}

// DiffOptions represents the CLI options for the diff command
type DiffOptions struct {
	Options
	Write bool
}

// AddFlags adds the flags for this set of options to the given command
func (opts *DiffOptions) AddFlags(cmd *cobra.Command) {
	opts.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&opts.Write, "write", "w", false, "Write patch files")
}

// CheckOptions represents the CLI options for the check command
type CheckOptions struct {
	Options
	Upstream bool
}

// AddFlags adds the flags for this set of options to the given command
func (opts *CheckOptions) AddFlags(cmd *cobra.Command) {
	opts.Options.AddFlags(cmd)
	cmd.Flags().BoolVarP(&opts.Upstream, "upstream", "r", false, "check against the upstream")
}
