package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ObjectChecksum(v any) (string, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(v); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(b.Bytes())), nil
}

func FileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("file does not exist %s", path)
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func DirChecksum(path string) (string, error) {
	var checksums []string
	err := filepath.WalkDir(path, func(file string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		sum, err := FileChecksum(file)
		if err == nil {
			checksums = append(checksums, sum)
		}
		return err
	})
	if err != nil {
		return "", err
	}

	return sumOfSums(checksums)
}

func FilesChecksum(baseDir string, files []string) (string, error) {
	var checksums []string
	for _, file := range files {
		sum, err := FileChecksum(filepath.Join(baseDir, file))
		if err != nil {
			return "", err
		}
		checksums = append(checksums, sum)
	}

	return sumOfSums(checksums)
}

func sumOfSums(checksums []string) (string, error) {
	if len(checksums) == 0 {
		return "", fmt.Errorf("empty set")
	}
	if len(checksums) == 1 {
		return checksums[0], nil
	}
	sort.Strings(checksums)
	output := strings.Join(checksums, "\n")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(output))), nil
}
