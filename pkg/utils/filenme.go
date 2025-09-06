package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"unicode"
)

func BumpFilename(filename string) string {
	dir := filepath.Dir(filename)
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	fname := filename[:len(filename)-len(ext)]

	var n int
	mul := 1
	l := len(fname)
	for i := l - 1; i >= 0; i-- {
		if !unicode.IsDigit(rune(fname[i])) {
			break
		}
		d, _ := strconv.Atoi(string(fname[i]))
		d *= mul
		mul *= 10
		n = d + n
		l--
	}

	n++
	return filepath.Join(dir, fmt.Sprintf("%s%d%s", fname[:l], n, ext))
}
