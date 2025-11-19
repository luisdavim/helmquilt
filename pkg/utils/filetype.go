package utils

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// IsTextFile determines if a file is likely a text file by scanning for
// non-printable characters.
func IsTextFile(file io.ReadSeeker) (isText bool, rerr error) {
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

	if IsBinaryData(buffer[:n]) {
		return false, nil
	}

	return true, nil
}

// IsBinaryData checks if the given data is likely to be binary by scanning for
// non-printable characters.
func IsBinaryData(data []byte) bool {
	n := 4096
	if l := len(data); l < n {
		n = l
	}
	if utf8.ValidString(string(data[:n])) {
		return false
	}

	// Iterate through the bytes to check for non-printable characters.
	for i := 0; i < n; {
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError {
			// A decoding error indicates a non-UTF-8 sequence, likely binary.
			return true
		}
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			// If it's not a printable character or a space, it's likely binary.
			// However, we allow for some common control characters like tabs and newlines.
			if r != '\n' && r != '\r' && r != '\t' {
				return true
			}
		}
		i += size
	}

	return false
}
