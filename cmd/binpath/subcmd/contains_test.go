package subcmd

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestContainsCmd(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin"}, pathSep)

	tests := []struct {
		name   string
		env    string
		arg    string
		absent bool // true => expect ErrNotPresent (non-zero exit)
	}{
		{"present -> success", base, "/usr/bin", false},
		{"absent -> non-zero exit", base, "/nope", true},
		{"empty PATH -> absent", "", "/usr/bin", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.env)
			out, err := execCmd(t, newContainsCmd(), tt.arg)

			if out != "" {
				t.Errorf("contains %q: expected no output, got %q", tt.arg, out)
			}
			if tt.absent {
				if !errors.Is(err, ErrNotPresent) {
					t.Fatalf("contains %q: expected ErrNotPresent, got %v", tt.arg, err)
				}
			} else if err != nil {
				t.Fatalf("contains %q: unexpected error: %v", tt.arg, err)
			}
		})
	}
}

func TestContainsCmd_dependencyError(t *testing.T) {
	orig, had := os.LookupEnv("PATH")
	_ = os.Unsetenv("PATH")
	t.Cleanup(func() {
		if had {
			_ = os.Setenv("PATH", orig)
		}
	})

	out, err := execCmd(t, newContainsCmd(), "/usr/bin")
	if err == nil {
		t.Fatal("contains with PATH unset: expected error, got nil")
	}
	if errors.Is(err, ErrNotPresent) {
		t.Fatalf("dependency failure should not be reported as not-present: %v", err)
	}
	if !strings.Contains(out, "Error:") {
		t.Errorf("expected an error message on stderr, got %q", out)
	}
}

func TestContainsCmd_requiresExactlyOneArg(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")
	if _, err := execCmd(t, newContainsCmd()); err == nil {
		t.Fatal("contains with no args: expected error, got nil")
	}
}
