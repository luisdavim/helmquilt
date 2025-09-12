package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
)

func diffCmd() *cobra.Command {
	var opts config.DiffOptions

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Check if changes were made and return the differences",
		Long: `
Check if changes were made and return the differences.
It will pull the charts into a temporary location using the config and compare that with the current state of the WorkDir.
By default it will only print the diffs but it can also store them in files and autommatically add them to the config.

Note that if you don't pin the source chart version and the remote chart is updaed, diff will return the differences between the local and the remote.
In this case, t is advisable that you run diff, or check with --upstream before maiking any changes to check if the upstream has changed.`,
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
			if !opts.Write {
				opts.DryRun = true
			}
			changed, err := helmquilt.Diff(cmd.Context(), opts)

			if len(changed) != 0 {
				fmt.Fprintln(os.Stderr, "Changes where detected on the following charts:")
				for _, name := range changed {
					fmt.Fprintf(os.Stderr, "\t-%s\n", name)
				}
			}

			return checkErr(err)
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
