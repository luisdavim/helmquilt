package utils

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/diff"
)

func DiffFiles(oldName, newName string) ([]byte, error) {
	oldB, err := os.ReadFile(oldName)
	if err != nil {
		return nil, err
	}
	newB, err := os.ReadFile(newName)
	if err != nil {
		return nil, err
	}
	diff := diff.Diff(oldName, oldB, newName, newB)

	return diff, nil
}

func DiffDirs(oldDir, newDir string) ([]byte, error) {
	var allDiffs bytes.Buffer

	oldFiles, err := readDir(oldDir)
	if err != nil {
		return nil, fmt.Errorf("faild to read from %s: %w", oldDir, err)
	}
	newFiles, err := readDir(newDir)
	if err != nil {
		return nil, fmt.Errorf("faild to read from %s: %w", newDir, err)
	}

	for oldName := range oldFiles {
		oldB, err := os.ReadFile(oldName)
		if err != nil {
			return nil, err
		}
		var newB []byte
		relName, err := filepath.Rel(oldDir, oldName)
		if err != nil {
			return nil, err
		}
		newName := filepath.Join(newDir, relName)
		if ok := newFiles[newName]; ok {
			newB, err = os.ReadFile(newName)
			if err != nil {
				return nil, err
			}
		}
		if newB == nil {
			newB = []byte("")
		}

		d := diff.Diff(relName, oldB, relName, newB)
		if len(d) != 0 {
			_, _ = allDiffs.WriteString("\n")
			_, _ = allDiffs.Write(d)
		}
		delete(newFiles, newName)
	}

	for newName := range newFiles {
		newB, err := os.ReadFile(newName)
		if err != nil {
			return nil, err
		}
		relName, err := filepath.Rel(newDir, newName)
		if err != nil {
			return nil, err
		}

		d := diff.Diff(relName, []byte(""), relName, newB)
		if len(d) != 0 {
			_, _ = allDiffs.WriteString("\n")
			_, _ = allDiffs.Write(d)
		}
	}

	if allDiffs.Len() > 0 {
		_, _ = allDiffs.WriteString("\n")
	}
	return allDiffs.Bytes(), nil
}

func readDir(root string) (map[string]bool, error) {
	files := make(map[string]bool)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files[path] = true
		return nil
	})

	return files, err
}
