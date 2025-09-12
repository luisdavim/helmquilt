package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func applyCmd() *cobra.Command {
	var opts config.ApplyOptions

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the helmquilt config, fetch charts and apply patches",
		Long: `
Apply the helmquilt config, fetch charts and apply patches.
In the end it will generate a lock file that is used on subsequent runs to determine what needs to be updated or not.`,
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
			return checkErr(helmquilt.Run(cmd.Context(), helmquilt.ApplyAction, opts))
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
