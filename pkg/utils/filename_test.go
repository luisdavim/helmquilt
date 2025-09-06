package utils

import "testing"

func TestBunmpFilename(t *testing.T) {
	tests := []struct {
		name     string
		in       string
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BumpFilename(tc.in)
			if got != tc.expected {
				t.Errorf("Unexpected result; wanted: %s, got: %s", tc.expected, got)
			}
		})
	}
}
