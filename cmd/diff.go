package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func diffCmd() *cobra.Command {
	var opts config.DiffOptions

	cmd := &cobra.Command{
		Use:          "diff",
		Short:        "check if changes were made and return the differences",
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
			return checkErr(helmquilt.Diff(cmd.Context(), opts))
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
