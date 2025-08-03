package main

import (
	"errors"
	"os"

	"github.com/luisdavim/helmquilt/cmd"
)

func main() {
	rootCmd := cmd.New()

	if err := rootCmd.Execute(); err != nil {
		// only print usage on usage errors
		if !errors.Is(err, cmd.ErrHelmquilt) {
			_ = rootCmd.Usage()
		}
		os.Exit(1)
	}
}
