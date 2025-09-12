package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func checkCmd() *cobra.Command {
	var (
		opts     config.Options
		upstream bool
	)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.ConfigFile, err = filepath.Abs(opts.ConfigFile)
			if err != nil {
				return checkErr(err)
			}
			if opts.WorkDir == "" {
				opts.WorkDir = filepath.Dir(opts.ConfigFile)
			}

			opts.DryRun = true

			if upstream {
				changed, err := helmquilt.Diff(cmd.Context(), config.DiffOptions{Options: opts})

				if len(changed) != 0 {
					fmt.Fprintln(os.Stderr, "\nChanges where detected on the following charts:")
					for _, name := range changed {
						fmt.Fprintf(os.Stderr, "\t-%s\n", name)
					}
					fmt.Fprintln(os.Stderr, "")

					if err == nil {
						err = helmquilt.ErrChartsChanged
					}
				}

				return checkErr(err)
			}
			return checkErr(helmquilt.Run(cmd.Context(), helmquilt.CheckAction, config.ApplyOptions{Options: opts}))
		},
	}

	opts.AddFlags(cmd)
	cmd.Flags().BoolVarP(&upstream, "upstream", "r", false, "check against the upstream")

	return cmd
}
