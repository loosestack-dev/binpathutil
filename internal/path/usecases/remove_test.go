package usecases_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"binpathutil/internal/path/usecases"
)

func TestRemove(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		element string
		path    string
		want    string
		wantErr bool
	}{
		{"present", "/usr/bin", base, "/bin", false},
		{"absent", "/opt/bin", base, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.Remove(tt.element, stubErr(tt.path, nil))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Remove(%q): expected error, got nil (result %q)", tt.element, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Remove(%q) unexpected error: %v", tt.element, err)
			}
			if got != tt.want {
				t.Errorf("Remove(%q) = %q, want %q", tt.element, got, tt.want)
			}
		})
	}
}

func TestRemove_dependencyError(t *testing.T) {
	got, err := usecases.Remove("/usr/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("Remove: expected error when dependency fails, got nil (result %q)", got)
	}
}

func TestRemoveIfPresent(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	tests := []struct {
		name    string
		element string
		path    string
		want    string
	}{
		{"present", "/usr/bin", base, "/bin"},
		{"absent (no-op)", "/opt/bin", base, base},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.RemoveIfPresent(tt.element, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("RemoveIfPresent(%q) unexpected error: %v", tt.element, err)
			}
			if got != tt.want {
				t.Errorf("RemoveIfPresent(%q) = %q, want %q", tt.element, got, tt.want)
			}
		})
	}
}

func TestRemoveIfPresent_dependencyError(t *testing.T) {
	got, err := usecases.RemoveIfPresent("/usr/bin", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("RemoveIfPresent: expected error when dependency fails, got nil (result %q)", got)
	}
}

func TestRemoveAllOccurence(t *testing.T) {
	sep := string(os.PathListSeparator)

	tests := []struct {
		name    string
		element string
		path    string
		want    string
	}{
		{"multiple", "/a", strings.Join([]string{"/a", "/x", "/a", "/a"}, sep), "/x"},
		{"single", "/a", strings.Join([]string{"/a", "/x"}, sep), "/x"},
		{"absent (no-op)", "/a", strings.Join([]string{"/x", "/y"}, sep), strings.Join([]string{"/x", "/y"}, sep)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.RemoveAllOccurence(tt.element, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("RemoveAllOccurence(%q) unexpected error: %v", tt.element, err)
			}
			if got != tt.want {
				t.Errorf("RemoveAllOccurence(%q) = %q, want %q", tt.element, got, tt.want)
			}
		})
	}
}

func TestRemoveAllOccurence_dependencyError(t *testing.T) {
	got, err := usecases.RemoveAllOccurence("/a", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("RemoveAllOccurence: expected error when dependency fails, got nil (result %q)", got)
	}
}
