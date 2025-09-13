package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func checkCmd() *cobra.Command {
	var opts config.CheckOptions

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Read the helmquilt config and check if everything is up-to-date with the lock file",
		Long: `
Read the helmquilt config and check if everything is up-to-date with the lock file.
It will exit with a non zero code if the charts do not match the lock file.

By default the check command only looks at the lock file, if you don't pin the source versions,
you can use --upstream to check against the upstream charts.`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.Quiet, err = getQuietOption(cmd)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.ConfigFile, err = filepath.Abs(opts.ConfigFile)
			if err != nil {
				return checkErr(err)
			}
			if opts.WorkDir == "" {
				opts.WorkDir = filepath.Dir(opts.ConfigFile)
			}

			return checkErr(helmquilt.Check(cmd.Context(), opts))
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
