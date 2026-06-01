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

func TestRemoveRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin", "/usr/local/bin"}, sep)

	tests := []struct {
		name    string
		pattern string
		path    string
		want    string
		wantErr bool
	}{
		{"removes first match", "usr", base, strings.Join([]string{"/bin", "/usr/local/bin"}, sep), false},
		{"anchored match", "^/bin$", base, strings.Join([]string{"/usr/bin", "/usr/local/bin"}, sep), false},
		{"no match errors (strict)", "^/opt", base, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.RemoveRegex(tt.pattern, stubErr(tt.path, nil))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("RemoveRegex(%q): expected error, got nil (result %q)", tt.pattern, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("RemoveRegex(%q) unexpected error: %v", tt.pattern, err)
			}
			if got != tt.want {
				t.Errorf("RemoveRegex(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestRemoveRegex_invalidRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	got, err := usecases.RemoveRegex("[", stubErr(base, nil))
	if err == nil {
		t.Fatalf("RemoveRegex: expected error for malformed pattern, got nil (result %q)", got)
	}
}

func TestRemoveRegex_dependencyError(t *testing.T) {
	got, err := usecases.RemoveRegex("usr", stubErr("", errors.New("PATH unavailable")))
	if err == nil {
		t.Fatalf("RemoveRegex: expected error when dependency fails, got nil (result %q)", got)
	}
}

func TestRemoveIfPresentRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin", "/usr/local/bin"}, sep)

	tests := []struct {
		name    string
		pattern string
		path    string
		want    string
	}{
		{"removes first match", "usr", base, strings.Join([]string{"/bin", "/usr/local/bin"}, sep)},
		{"no match is no-op", "^/opt", base, base},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.RemoveIfPresentRegex(tt.pattern, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("RemoveIfPresentRegex(%q) unexpected error: %v", tt.pattern, err)
			}
			if got != tt.want {
				t.Errorf("RemoveIfPresentRegex(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestRemoveIfPresentRegex_invalidRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	if _, err := usecases.RemoveIfPresentRegex("[", stubErr(base, nil)); err == nil {
		t.Fatal("RemoveIfPresentRegex: expected error for malformed pattern, got nil")
	}
}

func TestRemoveIfPresentRegex_dependencyError(t *testing.T) {
	if _, err := usecases.RemoveIfPresentRegex("usr", stubErr("", errors.New("PATH unavailable"))); err == nil {
		t.Fatal("RemoveIfPresentRegex: expected error when dependency fails, got nil")
	}
}

func TestRemoveAllOccurenceRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin", "/usr/local/bin"}, sep)

	tests := []struct {
		name    string
		pattern string
		path    string
		want    string
	}{
		{"removes all matches", "usr", base, "/bin"},
		{"removes single match", "^/bin$", base, strings.Join([]string{"/usr/bin", "/usr/local/bin"}, sep)},
		{"no match is no-op", "^/opt", base, base},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.RemoveAllOccurenceRegex(tt.pattern, stubErr(tt.path, nil))
			if err != nil {
				t.Fatalf("RemoveAllOccurenceRegex(%q) unexpected error: %v", tt.pattern, err)
			}
			if got != tt.want {
				t.Errorf("RemoveAllOccurenceRegex(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestRemoveAllOccurenceRegex_invalidRegex(t *testing.T) {
	sep := string(os.PathListSeparator)
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep)

	if _, err := usecases.RemoveAllOccurenceRegex("[", stubErr(base, nil)); err == nil {
		t.Fatal("RemoveAllOccurenceRegex: expected error for malformed pattern, got nil")
	}
}

func TestRemoveAllOccurenceRegex_dependencyError(t *testing.T) {
	if _, err := usecases.RemoveAllOccurenceRegex("usr", stubErr("", errors.New("PATH unavailable"))); err == nil {
		t.Fatal("RemoveAllOccurenceRegex: expected error when dependency fails, got nil")
	}
}
