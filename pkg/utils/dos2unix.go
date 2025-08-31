package utils

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Dos2UnixDir(path string) error {
	return filepath.WalkDir(path, func(file string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return Dos2Unix(file)
	})
}

func Dos2Unix(filename string) error {
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	// Write to a temp file
	tmpfile, err := os.CreateTemp("", "dos2unix")
	if err != nil {
		return err
	}
	defer func() { _ = tmpfile.Close() }()

	reader := bufio.NewReader(infile)
	writer := bufio.NewWriter(tmpfile)

	for {
		for isPrefix := true; isPrefix; {
			var buf []byte
			buf, isPrefix, err = reader.ReadLine()
			if err != nil && err != io.EOF {
				return err
			}
			_, _ = writer.Write(buf)
		}
		_, _ = writer.WriteString("\n")
		if err == io.EOF {
			break
		}
	}
	_ = writer.Flush()

	// Replace the original file with the new one
	if err := os.Rename(tmpfile.Name(), filename); err != nil {
		return err
	}

	return nil
}
