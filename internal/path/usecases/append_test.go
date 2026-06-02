package usecases_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/loosestack-dev/binpathutil/internal/path/usecases"
)

func stubErr(path string, err error) func() (string, error) {
	return func() (string, error) { return path, err }
}

func TestAppend(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	got, err := usecases.Append("/opt/bin", stubErr(base, nil))
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}
	want := base + sep + "/opt/bin"
	if got != want {
		t.Errorf("Append(%q) = %q, want %q", "/opt/bin", got, want)
	}
}

func TestAppendIfAbsent(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		element string
		path    string
		want    string
	}{
		{"already present", "/usr/bin", base, base},
		{"absent", "/opt/bin", base, base + sep + "/opt/bin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.AppendIfAbsent(tt.element, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("AppendIfAbsent returned unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("AppendIfAbsent(%q) = %q, want %q", tt.element, got, tt.want)
			}
		})
	}
}

func TestAppend_dependencyError(t *testing.T) {
	got, err := usecases.Append("/opt/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("Append: expected error when dependency fails, got nil (result %q)", got)
	}
}
