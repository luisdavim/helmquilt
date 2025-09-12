package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// BumpFilename looks for a number at the end of the given file name and incremetns it by one
// if the file name has no numbers at the end 1 is added to it
// if a delimiter is provided, the number is added afer it
// if a delimiter is provided and the given filename doesn't have it any number at the end of the file is ignored
// and the delimiter followed by 1 is added to the end of the file
func BumpFilename(filename, delim string, width int) string {
	// decompose name
	dir := filepath.Dir(filename)
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	fname := filename[:len(filename)-len(ext)]

	// handle the delimiter and only take into account the parts of the name after the last occurence of the delimiter
	var prefix string
	if delim != "" {
		idx := strings.LastIndex(fname, delim)
		if idx > 0 {
			prefix = fname[:idx]
			fname = fname[idx+1:]
		} else {
			prefix = fname
			fname = ""
		}
		prefix += delim
	}

	// look for a number at the end of the file name
	var n int
	mul := 1
	l := len(fname)
	for i := l - 1; i >= 0; i-- {
		if !unicode.IsDigit(rune(fname[i])) {
			break
		}
		// we're reading the filename backwards so we have to reverse the number
		d, _ := strconv.Atoi(string(fname[i]))
		n = d*mul + n
		mul *= 10
		// subtract the number length from the filename length so we can replace it
		l--
	}

	// bump the file number
	n++
	// recompose name
	return filepath.Join(dir, fmt.Sprintf("%s%s%0*d%s", prefix, fname[:l], width, n, ext))
}
