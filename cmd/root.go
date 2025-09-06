package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/luisdavim/helmquilt/pkg/config"
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
		quiet bool
		out   io.Writer
	)

	rootCmd := &cobra.Command{
		Use:          "helmquilt <apply|check|diff>",
		Short:        "helmquilt is a tool for managing helm package patches",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			out = os.Stderr
			if quiet {
				out = io.Discard
			}
			ctx := logger.NewContext(cmd.Context(), "helmquilt", out)
			cmd.SetContext(ctx)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence the logs")
	rootCmd.AddCommand(applyCmd(), checkCmd(), diffCmd())

	return rootCmd
}
