package helmquilt

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bluekeyes/go-gitdiff/gitdiff"

	"github.com/luisdavim/helmquilt/pkg/config"
	"github.com/luisdavim/helmquilt/pkg/logger"
	"github.com/luisdavim/helmquilt/pkg/utils"
)

func keepOnly(ctx context.Context, destDir string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	logger := logger.FromContext(ctx)
	keep := make(map[string]bool)

	for _, filePath := range files {
		fullPath := filepath.Join(destDir, filePath)
		if _, err := os.Stat(fullPath); err == nil {
			keep[fullPath] = true
		} else {
			return fmt.Errorf("keep file: %w", err)
		}
	}

	if err := filepath.Walk(destDir, func(path string, info fs.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			if !keep[path] {
				logger.Println("Removing", path)
				if err := os.RemoveAll(path); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// clean up left over empty directories
	return utils.RemoveEmptyFolders(destDir)
}

func applyPatches(ctx context.Context, chart config.Chart, workDir, destDir string) error {
	logger := logger.FromContext(ctx)

	for _, patchFile := range chart.Patches {
		patchFile = filepath.Join(workDir, config.PatchesPath, patchFile)
		logger.Printf("Applying patch %q to %q\n", patchFile, chart.GetName())

		patch, err := os.Open(patchFile)
		if err != nil {
			return fmt.Errorf("failed to open patch file: %w", err)
		}

		// files is alist of files referenced in the patch
		files, _, err := gitdiff.Parse(patch)
		if err != nil {
			return err
		}

		for _, f := range files {
			f.OldName = filepath.Join(destDir, f.OldName)
			f.NewName = filepath.Join(destDir, f.NewName)
			logger.Printf("Patching file: %q\n", f.OldName)

			// apply the changes in the patch to a source file
			src, err := os.Open(f.OldName)
			if err != nil {
				return fmt.Errorf("failed to open file to be patched: %w", err)
			}

			tmpfile, err := os.CreateTemp("", "helmquilt")
			if err != nil {
				_ = src.Close()
				return fmt.Errorf("failed to create temp file: %w", err)
			}

			if err := gitdiff.Apply(tmpfile, src, f); err != nil {
				_ = src.Close()
				_ = tmpfile.Close()
				return fmt.Errorf("failed to apply patch: %w", err)
			}

			if err := os.Remove(src.Name()); err != nil {
				return err
			}
			err = os.Rename(tmpfile.Name(), f.NewName)
			_ = src.Close()
			_ = tmpfile.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func applyFileMigrations(ctx context.Context, chart config.Chart, workDir, destDir string) error {
	logger := logger.FromContext(ctx)

	if len(chart.Keep) > 0 && len(chart.Remove) > 0 {
		return fmt.Errorf("keep and remove can't be used together")
	}

	for _, move := range chart.Move {
		logger.Printf("Renaming %q to %q\n", move.Source, move.Dest)
		if err := utils.MoveFile(destDir, move.Source, move.Dest); err != nil {
			return err
		}
	}

	for _, file := range chart.Remove {
		logger.Printf("Removing %s\n", file)
		if err := os.RemoveAll(filepath.Join(destDir, file)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if err := keepOnly(ctx, destDir, chart.Keep); err != nil {
		return err
	}

	return applyPatches(ctx, chart, workDir, destDir)
}
