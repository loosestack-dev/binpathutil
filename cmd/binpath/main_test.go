package main

import (
	"errors"
	"fmt"
	"testing"

	"binpathutil/cmd/binpath/subcmd"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"success", nil, 0},
		{"absent (sentinel)", subcmd.ErrNotPresent, 1},
		{"absent (wrapped)", fmt.Errorf("contains failed: %w", subcmd.ErrNotPresent), 1},
		{"genuine error", errors.New("PATH unavailable"), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.err); got != tt.want {
				t.Errorf("exitCode(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}
