package usecases_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"binpathutil/internal/path/usecases"
)

func TestPrepend(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	got, err := usecases.Prepend("/opt/bin", stubErr(base, nil))
	if err != nil {
		t.Fatalf("Prepend returned unexpected error: %v", err)
	}
	want := "/opt/bin" + sep + base
	if got != want {
		t.Errorf("Prepend(%q) = %q, want %q", "/opt/bin", got, want)
	}
}

func TestPrependIfAbsent(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		element string
		path    string
		want    string
	}{
		{"already present", "/usr/bin", base, base},
		{"absent", "/opt/bin", base, "/opt/bin" + sep + base},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.PrependIfAbsent(tt.element, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("PrependIfAbsent returned unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("PrependIfAbsent(%q) = %q, want %q", tt.element, got, tt.want)
			}
		})
	}
}

func TestPrepend_dependencyError(t *testing.T) {
	got, err := usecases.Prepend("/opt/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("Prepend: expected error when dependency fails, got nil (result %q)", got)
	}
}
