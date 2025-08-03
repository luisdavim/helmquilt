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

type LockFile struct {
	Charts []ChartLock `json:"charts"`
	Hash   string      `json:"hash"`
}

type ChartLock struct {
	Name       string          `json:"name"`
	Version    string          `json:"version"`
	Components []CharPathtLock `json:"components"`
}

type CharPathtLock struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

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

func UpdateLockfile(cfg Config, opts Options) error {
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

	if err := os.WriteFile(filepath.Join(opts.WorkDir, lockFilename), out, 0o644); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}

	// TODO: is this needed? in the current implementaion this will drop any file comments
	return Save(cfg, opts)
}
