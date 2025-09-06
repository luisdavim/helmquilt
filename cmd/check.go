package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func checkCmd() *cobra.Command {
	var opts config.Options

	cmd := &cobra.Command{
		Use:          "check",
		Short:        "read the helmquilt config and lock files and check if everything is up-to-date",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.ConfigFile, err = filepath.Abs(opts.ConfigFile)
			if err != nil {
				return checkErr(err)
			}
			if opts.WorkDir == "" {
				opts.WorkDir = filepath.Dir(opts.ConfigFile)
			}
			return checkErr(helmquilt.Run(cmd.Context(), helmquilt.CheckAction, opts))
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
