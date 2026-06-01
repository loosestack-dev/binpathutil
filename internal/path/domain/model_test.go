package domain_test

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"binpathutil/internal/path/domain"
)

func sep() string { return string(os.PathListSeparator) }

func TestNewEnvPath(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{"empty", "", []string{}},
		{"single", "/usr/bin", []string{"/usr/bin"}},
		{"multiple", strings.Join([]string{"/usr/bin", "/bin"}, sep()), []string{"/usr/bin", "/bin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.NewEnvPath(tt.raw).Entries
			if !slices.Equal(got, tt.want) {
				t.Errorf("NewEnvPath(%q).Entries = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestEnvPath_Contains(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		entry string
		want  bool
	}{
		{"present", strings.Join([]string{"/usr/bin", "/bin"}, sep()), "/bin", true},
		{"absent", strings.Join([]string{"/usr/bin", "/bin"}, sep()), "/opt/bin", false},
		{"empty path", "", "/bin", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domain.NewEnvPath(tt.raw).Contains(tt.entry); got != tt.want {
				t.Errorf("Contains(%q) = %v, want %v", tt.entry, got, tt.want)
			}
		})
	}
}

func TestEnvPath_ContainsMatch(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep())

	tests := []struct {
		name    string
		raw     string
		pattern string
		want    bool
	}{
		{"matches one entry", base, "bin$", true},
		{"substring match", base, "usr", true},
		{"matches none", base, "^/opt", false},
		{"empty path", "", "bin", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := regexp.MustCompile(tt.pattern)
			if got := domain.NewEnvPath(tt.raw).ContainsMatch(re); got != tt.want {
				t.Errorf("ContainsMatch(%q) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestEnvPath_Append(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		entry string
		want  []string
	}{
		{"to non-empty", "/usr/bin", "/opt/bin", []string{"/usr/bin", "/opt/bin"}},
		{"to empty", "", "/opt/bin", []string{"/opt/bin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewEnvPath(tt.raw)
			p.Append(tt.entry)
			if !slices.Equal(p.Entries, tt.want) {
				t.Errorf("after Append(%q), Entries = %v, want %v", tt.entry, p.Entries, tt.want)
			}
		})
	}
}

func TestEnvPath_Prepend(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		entry string
		want  []string
	}{
		{"to non-empty", "/usr/bin", "/opt/bin", []string{"/opt/bin", "/usr/bin"}},
		{"to empty", "", "/opt/bin", []string{"/opt/bin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewEnvPath(tt.raw)
			p.Prepend(tt.entry)
			if !slices.Equal(p.Entries, tt.want) {
				t.Errorf("after Prepend(%q), Entries = %v, want %v", tt.entry, p.Entries, tt.want)
			}
		})
	}
}

func TestEnvPath_Remove(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin"}, sep())

	tests := []struct {
		name    string
		raw     string
		entry   string
		want    []string
		wantErr bool
	}{
		{"existing", base, "/usr/bin", []string{"/bin"}, false},
		{"missing", base, "/opt/bin", []string{"/usr/bin", "/bin"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewEnvPath(tt.raw)
			err := p.Remove(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove(%q) error = %v, wantErr %v", tt.entry, err, tt.wantErr)
			}
			if !slices.Equal(p.Entries, tt.want) {
				t.Errorf("after Remove(%q), Entries = %v, want %v", tt.entry, p.Entries, tt.want)
			}
		})
	}
}

func TestEnvPath_RemoveMatch(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin", "/usr/local/bin"}, sep())

	tests := []struct {
		name    string
		raw     string
		pattern string
		want    []string
		wantErr bool
	}{
		{"removes first match", base, "usr", []string{"/bin", "/usr/local/bin"}, false},
		{"anchored single match", base, "^/bin$", []string{"/usr/bin", "/usr/local/bin"}, false},
		{"no match errors", base, "^/opt", []string{"/usr/bin", "/bin", "/usr/local/bin"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewEnvPath(tt.raw)
			err := p.RemoveMatch(regexp.MustCompile(tt.pattern))
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveMatch(%q) error = %v, wantErr %v", tt.pattern, err, tt.wantErr)
			}
			if !slices.Equal(p.Entries, tt.want) {
				t.Errorf("after RemoveMatch(%q), Entries = %v, want %v", tt.pattern, p.Entries, tt.want)
			}
		})
	}
}

func TestEnvPath_String(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"empty", "", ""},
		{"single", "/usr/bin", "/usr/bin"},
		{"multiple", strings.Join([]string{"/usr/bin", "/bin"}, sep()), strings.Join([]string{"/usr/bin", "/bin"}, sep())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domain.NewEnvPath(tt.raw).String(); got != tt.want {
				t.Errorf("NewEnvPath(%q).String() = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
