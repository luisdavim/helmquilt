package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
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

	if errors.Is(err, config.ErrMissingConfig) {
		return err
	}

	return fmt.Errorf("%w%w", ErrHelmquilt, err)
}

func New() *cobra.Command {
	var (
		opts  config.Options
		quiet bool
		out   io.Writer
	)

	rootCmd := &cobra.Command{
		Use:          "helmquilt <apply|check|diff>",
		Short:        "helmquilt is a tool for managing helm package patches",
		Args:         cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs:    []cobra.Completion{"apply", "check", "diff"},
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
			out = os.Stderr
			if quiet {
				out = io.Discard
			}
			ctx := logger.NewContext(cmd.Context(), "helmquilt", out)
			return checkErr(helmquilt.Run(ctx, args[0], opts))
		},
	}

	opts.AddFlags(rootCmd)
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Silence the logs")

	return rootCmd
}
