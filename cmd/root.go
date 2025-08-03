package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/helmquilt"
	"github.com/luisdavim/helmquilt/pkg/logger"
)

var ErrHelmquilt = errors.New("")

// wrap the returned error with a custom one so we can distinguish usage errors and print the usage help
func checkErr(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w%w", ErrHelmquilt, err)
}

func New() *cobra.Command {
	var opts config.Options

	rootCmd := &cobra.Command{
		Use:          "helmquilt <apply|check>",
		Short:        "helmquilt is a tool for managing helm package patches",
		Args:         cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs:    []cobra.Completion{"apply", "check"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.ConfigFile, err = filepath.Abs(opts.ConfigFile)
			if err != nil {
				return checkErr(err)
			}
			opts.WorkDir = filepath.Dir(opts.ConfigFile)
			ctx := logger.NewContext(cmd.Context(), "helmquilt")
			return checkErr(helmquilt.Run(ctx, args[0], opts))
		},
	}
	rootCmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force run (ignore lock file)")
	rootCmd.Flags().BoolVarP(&opts.Repack, "repack", "r", false, "Repack the chart as a tarball")
	rootCmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "./helmquilt.yaml", "path to the config file")

	return rootCmd
}
