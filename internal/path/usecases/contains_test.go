package usecases_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"binpathutil/internal/path/usecases"
)

func TestContains(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		element string
		path    string
		want    bool
	}{
		{"present", "/usr/bin", base, true},
		{"absent", "/opt/bin", base, false},
		{"empty path", "/bin", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.Contains(tt.element, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("Contains(%q) unexpected error: %v", tt.element, err)
			}
			if got != tt.want {
				t.Errorf("Contains(%q) = %v, want %v", tt.element, got, tt.want)
			}
		})
	}
}

// When the dependency fails, the error is propagated and the result is false.
func TestContains_dependencyError(t *testing.T) {
	got, err := usecases.Contains("/usr/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("Contains: expected error when dependency fails, got nil")
	}
	if got {
		t.Errorf("Contains on dependency error = %v, want false", got)
	}
}
