package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
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

	// TODO: the following 2 checks could be optimised by reading the first 4Kb of the file once and reuseing it for each check
	if isText, err := isTextFile(infile); !isText || err != nil {
		return err
	}

	if hasCR, err := hasDOSLineEndings(infile); !hasCR || err != nil {
		return err
	}

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
		if err == io.EOF {
			break
		}
		_, _ = writer.WriteString("\n")
	}
	_ = writer.Flush()

	// Replace the original file with the new one
	if err := os.Rename(tmpfile.Name(), filename); err != nil {
		return err
	}

	return nil
}

func hasDOSLineEndings(file io.ReadSeeker) (foundCR bool, rerr error) {
	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // Read in 4KB chunks

	defer func() {
		if rerr != nil {
			return
		}
		// move the reader back to the start
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			rerr = err
		}
	}()

	for {
		// Read a chunk of the file
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, err
		}

		// Iterate through the bytes in the buffer
		for i := range n {
			switch buffer[i] {
			case '\r':
				// We found a carriage return, set the flag
				foundCR = true
			case '\n':
				if foundCR {
					// We found a newline immediately after a carriage return
					return true, nil
				}
				return false, nil
			default:
				// Reset the flag if the sequence is broken
				foundCR = false
			}
		}
	}

	return false, nil
}

// isTextFile determines if a file is likely a text file by scanning for
// non-printable characters.
func isTextFile(file io.ReadSeeker) (isText bool, rerr error) {
	defer func() {
		if rerr != nil {
			return
		}
		// move the reader back to the start
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			rerr = err
		}
	}()

	// Read the first 4Kb. This is usually sufficient.
	buffer := make([]byte, 4096)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("could not read file: %w", err)
	}

	if !utf8.ValidString(string(buffer[:n])) {
		return false, nil
	}

	// Iterate through the bytes to check for non-printable characters.
	for i := 0; i < n; {
		r, size := utf8.DecodeRune(buffer[i:])
		if r == utf8.RuneError {
			// A decoding error indicates a non-UTF-8 sequence, likely binary.
			return false, nil
		}
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			// If it's not a printable character or a space, it's likely binary.
			// However, we allow for some common control characters like tabs and newlines.
			if r != '\n' && r != '\r' && r != '\t' {
				return false, nil
			}
		}
		i += size
	}

	return true, nil
}
