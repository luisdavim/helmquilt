package utils

import (
	"testing"
)

func TestBunmpFilename(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		width    int
		delim    string
		expected string
	}{
		{
			name:     "no number",
			in:       "test.patch",
			expected: "test1.patch",
		},
		{
			name:     "zero",
			in:       "test0.patch",
			expected: "test1.patch",
		},
		{
			name:     "nine",
			in:       "test9.patch",
			expected: "test10.patch",
		},
		{
			name:     "nines",
			in:       "test99.patch",
			expected: "test100.patch",
		},
		{
			name:     "dash",
			in:       "test2-99.patch",
			expected: "test2-100.patch",
		},
		{
			name:     "with width",
			in:       "test9.patch",
			width:    3,
			expected: "test010.patch",
		},
		{
			name:     "with missing delim",
			in:       "test9.patch",
			delim:    "-",
			expected: "test9-1.patch",
		},
		{
			name:     "with delim",
			in:       "test-9.patch",
			delim:    "-",
			expected: "test-10.patch",
		},
		{
			name:     "with delim 2",
			in:       "test1-9.patch",
			delim:    "-",
			expected: "test1-10.patch",
		},
		{
			name:     "with delim 3",
			in:       "test1-a9.patch",
			delim:    "-",
			expected: "test1-a10.patch",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BumpFilename(tc.in, tc.delim, tc.width)
			if got != tc.expected {
				t.Errorf("Unexpected result; wanted: %s, got: %s", tc.expected, got)
			}
		})
	}
}
