package helmquilt

import (
	"context"
	"os"
	"path/filepath"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/logger"
	"github.com/luisdavim/helmquilt/pkg/utils"
)

// clearRemoved deletes charts found on the file system (and lock file) that are no longer in the configuration
func clearRemoved(ctx context.Context, cfg config.Config, lockfile *config.LockFile, opts config.Options) (bool, error) {
	logger := logger.FromContext(ctx)

	var (
		keep    []config.ChartLock
		changed bool
	)
	for _, lock := range lockfile.Charts {
		var found bool
		for _, chart := range cfg.Charts {
			if chart.Name == lock.Name {
				found = true
				break
			}
		}
		if found {
			keep = append(keep, lock)
			continue
		}
		// remove chart files
		for _, file := range lock.Components {
			if file.Path == "" || file.Path == config.PatchesPath || file.Path == filepath.Base(opts.ConfigFile) {
				// keep the patches and config file
				continue
			}
			changed = true
			if !opts.DryRun {
				logger.Printf("Cleaning up removed chart: %q\n", file.Path)
				_ = os.RemoveAll(filepath.Join(opts.WorkDir, file.Path))
			} else {
				logger.Printf("Chart: %q was removed from the config\n", file.Path)
			}
		}
	}
	lockfile.Charts = keep

	return changed, nil
}

// filterCurrent returns the subset of charts in the configuration that need to be updated
func filterCurrent(ctx context.Context, cfg config.Config, opts config.Options) (config.Config, bool, error) {
	logger := logger.FromContext(ctx)

	var changed bool
	if opts.Force {
		logger.Println("Force flag provided, updating")
		return cfg, changed, nil
	}

	lockfile, err := config.ReadLockFile(opts.WorkDir)
	if err != nil {
		logger.Println("No lockfile found, update required")
		return cfg, changed, nil
	}

	// delete charts referenced in the lockfile that are not in the new config
	changed, err = clearRemoved(ctx, cfg, &lockfile, opts)
	if err != nil {
		return cfg, changed, err
	}

	if len(cfg.Charts) > len(lockfile.Charts) {
		logger.Println("New charts where added to the config")
	}

	var modified config.Config
	for i, chart := range cfg.Charts {
		if i >= len(lockfile.Charts) {
			// new chart
			modified.Charts = append(modified.Charts, chart)
			continue
		}
		lock := lockfile.Charts[i]
		if lock.Name != chart.Name {
			logger.Println("Lockfile mismatch")
			modified = cfg
			break
		}
		if chart.Version != "" && lock.Version != chart.Version {
			logger.Printf("Chart %q version mismatch\n", chart.Name)
			modified.Charts = append(modified.Charts, chart)
			continue
		}

		for _, c := range lock.Components {
			var newSha string
			switch c.Path {
			case filepath.Base(opts.ConfigFile):
				newSha, err = utils.FileChecksum(opts.ConfigFile)
				if err != nil {
					return config.Config{}, changed, err
				}
			case config.PatchesPath:
				if len(chart.Patches) > 0 {
					newSha, err = utils.FilesChecksum(filepath.Join(opts.WorkDir, config.PatchesPath), chart.Patches)
					if err != nil {
						return config.Config{}, changed, err
					}
				}
			case config.OperationsPath:
				newSha, err = utils.ObjectChecksum(chart.FileOperations)
				if err != nil {
					return config.Config{}, changed, err
				}
			default:
				newSha, err = utils.DirChecksum(chart.GetFullName(opts.WorkDir))
				if err != nil {
					return config.Config{}, changed, err
				}
			}
			if newSha != c.Hash {
				logger.Printf("%s:%s hash mismatch; new: %q, old: %q\n", chart.Name, c.Path, newSha, c.Hash)
				modified.Charts = append(modified.Charts, chart)
				break
			}
		}
	}

	if len(modified.Charts) == 0 {
		// fallback for undetected changes
		if newSha, err := utils.FileChecksum(opts.ConfigFile); err != nil || lockfile.Hash != newSha {
			logger.Println("Config file hash mismatch: (uncaught changes)")
			return cfg, true, nil
		}
	} else {
		changed = true
	}

	return modified, changed, nil
}
