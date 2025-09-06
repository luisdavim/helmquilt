package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func BumpFilename(filename, delim string, width int) string {
	dir := filepath.Dir(filename)
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	fname := filename[:len(filename)-len(ext)]

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

	var n int
	mul := 1
	l := len(fname)
	for i := l - 1; i >= 0; i-- {
		if !unicode.IsDigit(rune(fname[i])) {
			break
		}
		d, _ := strconv.Atoi(string(fname[i]))
		n = d*mul + n
		mul *= 10
		l--
	}

	n++
	return filepath.Join(dir, fmt.Sprintf("%s%s%0*d%s", prefix, fname[:l], width, n, ext))
}
