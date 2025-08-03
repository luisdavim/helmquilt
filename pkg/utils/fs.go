package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveEmptyFolders(path string) error {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil
	}
	files, _ := os.ReadDir(path)
	for _, f := range files {
		fullpath := filepath.Join(path, f.Name())
		if stat, err := os.Stat(fullpath); err == nil && stat.IsDir() {
			RemoveEmptyFolders(fullpath)
		}
	}
	files, _ = os.ReadDir(path)
	if len(files) == 0 {
		os.Remove(path)
	}
	return nil
}

func MoveFile(destDir, source, dest string) error {
	sourceFile := filepath.Join(destDir, source)
	destFile := filepath.Join(destDir, dest)
	if _, err := os.Stat(sourceFile); err != nil {
		return fmt.Errorf("missing source file: %s", source)
	}
	if _, err := os.Stat(destFile); err == nil {
		return fmt.Errorf("dest file already exists: %s", dest)
	}
	if err := os.MkdirAll(filepath.Dir(destFile), 0o755); err != nil {
		return err
	}
	if err := os.Rename(sourceFile, destFile); err != nil {
		return err
	}
	// clean up left over empty directories
	return RemoveEmptyFolders(destDir)
}

func CopyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, _ := os.ReadFile(srcPath)
			os.WriteFile(dstPath, data, 0o644)
		}
	}
	return nil
}
