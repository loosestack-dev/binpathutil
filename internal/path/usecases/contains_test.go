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

func TestContains_dependencyError(t *testing.T) {
	got, err := usecases.Contains("/usr/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("Contains: expected error when dependency fails, got nil")
	}
	if got {
		t.Errorf("Contains on dependency error = %v, want false", got)
	}
}

func TestContainsRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{"substring match", "usr", base, true},
		{"anchored match", "^/usr", base, true},
		{"no match", "^/opt", base, false},
		{"empty path", "bin", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.ContainsRegex(tt.pattern, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("ContainsRegex(%q) unexpected error: %v", tt.pattern, err)
			}
			if got != tt.want {
				t.Errorf("ContainsRegex(%q) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestContainsRegex_invalidRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	got, err := usecases.ContainsRegex("[", stubErr(base, nil))
	if err == nil {
		t.Fatalf("ContainsRegex: expected error for malformed pattern, got nil")
	}
	if got {
		t.Errorf("ContainsRegex on invalid regex = %v, want false", got)
	}
}

func TestContainsRegex_dependencyError(t *testing.T) {
	got, err := usecases.ContainsRegex("usr", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("ContainsRegex: expected error when dependency fails, got nil")
	}
	if got {
		t.Errorf("ContainsRegex on dependency error = %v, want false", got)
	}
}
