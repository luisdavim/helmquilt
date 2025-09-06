package helmquilt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/logger"
	"github.com/luisdavim/helmquilt/pkg/utils"
)

func Run(ctx context.Context, action Action, opts config.Options) error {
	logger := logger.FromContext(ctx)

	if action == DiffAction {
		return Diff(ctx, config.DiffOptions{Options: opts})
	}

	cfg, err := config.Read(opts.ConfigFile)
	if err != nil {
		return err
	}

	if action == CheckAction {
		opts.DryRun = true
	}

	// skip charts that are up-to-date
	needUpdate, changed, err := filterCurrent(ctx, cfg, opts)
	if err != nil {
		return err
	}

	if len(needUpdate.Charts) == 0 {
		logger.Println("All is up to date!")
		return nil
	}

	if action == CheckAction {
		if changed {
			return fmt.Errorf("some charts need to be cleaned")
		}
		return fmt.Errorf("not all charts are up to date")
	}

	logger.Println("Updating Chart(s)")
	tempDir, _ := os.MkdirTemp("", "helmquilt")
	defer func() { _ = os.RemoveAll(tempDir) }()

	for _, chart := range needUpdate.Charts {
		destDir := filepath.Join(opts.WorkDir, chart.Path)
		chartDownloadPath, err := fetchChart(ctx, chart.Source, tempDir)
		if err != nil {
			return err
		}

		chartDestDir := filepath.Join(destDir, chart.Source.ChartName)
		if _, err := os.Stat(chartDestDir); err == nil {
			if err := os.RemoveAll(chartDestDir); err != nil {
				return err
			}
		}

		logger.Printf("Copying chart to %s\n", chartDestDir)
		if err := utils.CopyDir(filepath.Join(chartDownloadPath, chart.Source.ChartPath, chart.Source.ChartName), chartDestDir); err != nil {
			return err
		}

		if err := applyFileMigrations(ctx, chart, opts.WorkDir, chartDestDir); err != nil {
			return err
		}
		if err := utils.UpdateChartMetadata(chart.Version, chartDestDir); err != nil {
			return err
		}
		if err := utils.RemoveEmptyFolders(chartDestDir); err != nil {
			return err
		}
		if err := utils.Dos2UnixDir(chartDestDir); err != nil {
			return err
		}

		if opts.Repack {
			if chart.Version == "" {
				// the package file name needs to include the chart version
				chart.Version, err = utils.GetChartVersion(chartDestDir)
				if err != nil {
					return fmt.Errorf("failed to get chart version: %w", err)
				}
				for i := 0; i < len(cfg.Charts); i++ {
					if cfg.Charts[i].Name == chart.Name {
						cfg.Charts[i] = chart
						break
					}
				}
			}
			dst := fmt.Sprintf("%s-%s.tgz", chartDestDir, chart.Version)
			sum, err := utils.Compress(chartDestDir, dst)
			if err != nil {
				return fmt.Errorf("failed to repack %q: %w", chart.GetName(), err)
			}
			if err := os.RemoveAll(chartDestDir); err != nil {
				return err
			}
			logger.Printf("Chart saved to %q; checksum: %q\n", dst, sum)
		}
	}

	logger.Println("Updating lockfile")

	return config.UpdateLockfile(cfg, opts)
}
