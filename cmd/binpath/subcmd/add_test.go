package subcmd

import (
	"strings"
	"testing"
)

func TestAddCmd(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin"}, pathSep)

	tests := []struct {
		name string
		env  string
		args []string
		want string
	}{
		{"default prepends", base, []string{"/opt/bin"}, "/opt/bin" + pathSep + base},
		{"--first prepends", base, []string{"--first", "/opt/bin"}, "/opt/bin" + pathSep + base},
		{"--last appends", base, []string{"--last", "/opt/bin"}, base + pathSep + "/opt/bin"},
		{"--if-absent, present is no-op", base, []string{"-i", "/usr/bin"}, base},
		{"--if-absent, absent is added", base, []string{"-i", "/opt/bin"}, "/opt/bin" + pathSep + base},
		{"--last --if-absent, absent appends", base, []string{"-l", "-i", "/opt/bin"}, base + pathSep + "/opt/bin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.env)
			got, err := execCmd(t, newAddCmd(), tt.args...)
			if err != nil {
				t.Fatalf("add %v: unexpected error: %v", tt.args, err)
			}
			if got != tt.want {
				t.Errorf("add %v = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestAddCmd_mutuallyExclusiveFlags(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")
	if _, err := execCmd(t, newAddCmd(), "--first", "--last", "/opt/bin"); err == nil {
		t.Fatal("add --first --last: expected error, got nil")
	}
}

func TestAddCmd_requiresExactlyOneArg(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")
	if _, err := execCmd(t, newAddCmd()); err == nil {
		t.Fatal("add with no args: expected error, got nil")
	}
}
