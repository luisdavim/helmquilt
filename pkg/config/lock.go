package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/luisdavim/helmquilt/pkg/utils"
)

const (
	lockFilename   = "helmquilt.lock.yaml"
	PatchesPath    = "patches"
	OperationsPath = "operations"
)

// LockFile is the struture for the lock LockFile
// where the tool keeps the checksums from the last update
type LockFile struct {
	Charts []ChartLock `json:"charts"`
	Hash   string      `json:"hash"`
}

// ChartLock represents the state of an individual chart
type ChartLock struct {
	Name       string          `json:"name"`
	Version    string          `json:"version"`
	Components []CharPathtLock `json:"components"`
}

// CharPathtLock holds the checksums for the chart config components
type CharPathtLock struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// ReadLockFile reads the a lock file from the given path
func ReadLockFile(workDir string) (LockFile, error) {
	var locks LockFile
	lockFile := filepath.Join(workDir, lockFilename)
	if _, err := os.Stat(lockFile); errors.Is(err, os.ErrNotExist) {
		return locks, err
	}
	lockData, _ := os.ReadFile(lockFile)
	if err := yaml.Unmarshal(lockData, &locks); err != nil {
		return locks, err
	}

	return locks, nil
}

// UpdateLockfile calculates all the checksums and crates or updates a lock file in the workDir
func UpdateLockfile(cfg Config, opts ApplyOptions) error {
	locks := LockFile{Charts: []ChartLock{}}

	for _, chart := range cfg.Charts {
		sha, err := utils.FileChecksum(opts.ConfigFile)
		if err != nil {
			return err
		}
		locks.Hash = sha

		files := []CharPathtLock{}

		if len(chart.Patches) > 0 {
			sha, err := utils.FilesChecksum(filepath.Join(opts.WorkDir, PatchesPath), chart.Patches)
			if err != nil {
				return err
			}
			files = append(files, CharPathtLock{
				Path: PatchesPath,
				Hash: sha,
			})
		}

		if chart.HasOperations() {
			sha, err := utils.ObjectChecksum(chart.FileOperations)
			if err != nil {
				return err
			}
			files = append(files, CharPathtLock{
				Path: OperationsPath,
				Hash: sha,
			})
		}

		chartName := chart.GetName()
		chartPath := chart.GetFullName(opts.WorkDir)

		if opts.Repack {
			chartName = chart.GetTarName()
			chartPath = filepath.Join(opts.WorkDir, chartName)
			sha, err = utils.FileChecksum(chartPath)
			if err != nil {
				return err
			}
		} else {
			sha, err = utils.DirChecksum(chartPath)
			if err != nil {
				return err
			}
		}

		files = append(files, CharPathtLock{
			Path: chartName,
			Hash: sha,
		})

		locks.Charts = append(locks.Charts, ChartLock{
			Name:       chart.Name,
			Version:    chart.Version,
			Components: files,
		})
	}

	out, err := yaml.Marshal(locks)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return nil
	}

	if err := os.WriteFile(filepath.Join(opts.WorkDir, lockFilename), out, 0o644); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}

	// TODO: is this needed? in the current implementaion this will drop any file comments
	return Save(cfg, opts.Options)
}
