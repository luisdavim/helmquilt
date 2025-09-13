package helmquilt

import (
	"context"
	"fmt"
	"os"

	"github.com/luisdavim/helmquilt/pkg/config"
)

func Check(ctx context.Context, opts config.CheckOptions) error {
	opts.DryRun = true

	if opts.Upstream {
		changed, err := Diff(ctx, config.DiffOptions{Options: opts.Options})

		if len(changed) != 0 {
			fmt.Fprintln(os.Stderr, "\nChanges where detected on the following charts:")
			for _, name := range changed {
				fmt.Fprintf(os.Stderr, "\t-%s\n", name)
			}
			fmt.Fprintln(os.Stderr, "")

			if err == nil {
				err = ErrChartsChanged
			}
		}

		return err
	}
	return Run(ctx, CheckAction, config.ApplyOptions{Options: opts.Options})
}
