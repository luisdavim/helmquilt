package helmquilt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/logger"
	"github.com/luisdavim/helmquilt/pkg/utils"
)

func Diff(ctx context.Context, opts config.DiffOptions) ([]string, error) {
	logger := logger.FromContext(ctx)

	cfg, err := config.Read(opts.ConfigFile)
	if err != nil {
		return nil, err
	}

	tempDir, _ := os.MkdirTemp("", "helmquilt")
	defer func() { _ = os.RemoveAll(tempDir) }()

	tmpOpts := config.ApplyOptions{Options: opts.Options}
	tmpOpts.WorkDir = tempDir
	tmpOpts.Force = true

	if err := utils.CopyDir(filepath.Join(opts.WorkDir, config.PatchesPath), filepath.Join(tempDir, config.PatchesPath)); err != nil {
		return nil, err
	}

	logger.Println("Preparing current state for comparison!")
	if err := Run(ctx, ApplyAction, tmpOpts); err != nil {
		return nil, err
	}
	logger.Println("Current state ready!")

	var (
		updated bool
		changed []string
	)

	for i, chart := range cfg.Charts {
		oldChart := filepath.Join(tempDir, chart.Path, chart.Source.ChartName)
		newChart := filepath.Join(opts.WorkDir, chart.Path, chart.Source.ChartName)

		logger.Printf("Comparing %s with %s", oldChart, newChart)
		diff, err := utils.DiffDirs(oldChart, newChart)
		if err != nil {
			return changed, err
		}

		if len(diff) == 0 {
			logger.Println("No changes")
			continue
		}

		changed = append(changed, chart.Name)

		if opts.Write {
			file, err := getLatestPatch(chart.Name, opts.WorkDir)
			if err != nil {
				return changed, err
			}
			logger.Printf("Writing path file: %s", file)
			if err := os.WriteFile(file, diff, 0o644); err != nil {
				return changed, fmt.Errorf("faild to write patch file: %w", err)
			}
			chart.Patches = append(chart.Patches, filepath.Base(file))
			cfg.Charts[i] = chart
			updated = true
		}

		_, _ = os.Stdout.Write(diff)
	}

	if updated && opts.Write {
		logger.Println("Updating config file with new patches")
		if err := config.Save(cfg, opts.Options); err != nil {
			return changed, err
		}
		logger.Println("Updating lockfile")
		return changed, config.UpdateLockfile(cfg, config.ApplyOptions{Options: opts.Options})
	}

	if len(changed) == 0 {
		logger.Println("All is up to date!")
	}

	return changed, nil
}

func getLatestPatch(name, workDir string) (string, error) {
	files, err := utils.FindFile(filepath.Join(workDir, config.PatchesPath), fmt.Sprintf("%s*.patch", name))
	if err != nil {
		return "", err
	}

	file := filepath.Join(workDir, config.PatchesPath, fmt.Sprintf("%s.patch", name))
	if l := len(files); l > 0 {
		slices.Sort(files)
		file = files[l-1]
	}

	return utils.BumpFilename(file, "", 0), nil
}
