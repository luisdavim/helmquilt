package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// check for path traversal and correct forward slashes
func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

// Compress archives the contents of the source path and stoes the tarball in the destinaton path
func Compress(src, dst string) (string, error) {
	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}

	defer func() { _ = f.Close() }()

	return compress(src, f)
}

// Extract a tarball from the source path into the destinaton path
func Extract(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	return extract(f, dst)
}

func compress(src string, buf io.Writer) (hash string, rErr error) {
	h := sha256.New()
	mw := io.MultiWriter(h, buf)
	// tar > gzip > buf
	zr := gzip.NewWriter(mw)
	tw := tar.NewWriter(zr)

	defer func() {
		if err := tw.Close(); err != nil {
			rErr = errors.Join(rErr, err)
		}
		if err := zr.Close(); err != nil {
			rErr = errors.Join(rErr, err)
		}

		hash = fmt.Sprintf("%x", h.Sum(nil))
	}()

	fi, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	mode := fi.Mode()

	// single file
	if mode.IsRegular() {
		// get header
		header, err := tar.FileInfoHeader(fi, src)
		if err != nil {
			return "", err
		}
		// write header
		if err := tw.WriteHeader(header); err != nil {
			return "", err
		}
		// get content
		data, err := os.Open(src)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return "", err
		}

		return "", nil
	}

	if mode.IsDir() {
		// walk through every file in the folder
		err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// generate tar header
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// must provide real name
			// (see https://golang.org/src/archive/tar/common.go?#L626)
			header.Name = filepath.ToSlash(file)

			// write header
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// if not a dir, write file content
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err := io.Copy(tw, data); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return "", err
		}

		return "", nil
	}

	return "", fmt.Errorf("error: file type not supported")
}

func extract(src io.ReadSeeker, dst string) error {
	testBytes := make([]byte, 2)
	_, err := src.Read(testBytes)
	if err != nil {
		return err
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// if the archive is compressed, wrap the reader
	var reader io.Reader = src
	if isGzipCompressed(testBytes) {
		// ungzip
		reader, err = gzip.NewReader(src)
		if err != nil {
			return err
		}
	}

	// untar
	tr := tar.NewReader(reader)
	return untar(dst, tr)
}

func untar(dst string, tr *tar.Reader) error {
	// uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		target := header.Name

		// validate name against path traversal
		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", target)
		}

		// add dst + re-format slashes according to system
		target = filepath.Join(dst, header.Name)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// check the type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0o755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			// ensure the target dir exist
			targetDir := filepath.Dir(target)
			if _, err := os.Stat(targetDir); err != nil {
				if err := os.MkdirAll(targetDir, 0o755); err != nil {
					return err
				}
			}
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			_ = fileToWrite.Close()
		}
	}

	return nil
}

func isGzipCompressed(data []byte) bool {
	headerSize := 2

	if len(data) < headerSize {
		return false
	}

	gzipHeaderMagicNumber := []byte{0x1f, 0x8b}

	return bytes.Equal(data[:headerSize], gzipHeaderMagicNumber)
}
