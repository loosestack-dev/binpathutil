package subcmd

import (
	"strings"
	"testing"
)

func TestRemoveCmd(t *testing.T) {
	base := strings.Join([]string{"/usr/bin", "/bin"}, pathSep)
	dupes := strings.Join([]string{"/a", "/x", "/a", "/a"}, pathSep)

	tests := []struct {
		name string
		env  string
		args []string
		want string
	}{
		{"default removes present", base, []string{"/usr/bin"}, "/bin"},
		{"--if-present, absent is no-op", base, []string{"-i", "/nope"}, base},
		{"--if-present, present is removed", base, []string{"-i", "/usr/bin"}, "/bin"},
		{"--all strips duplicates", dupes, []string{"--all", "/a"}, "/x"},
		{"--all, absent is no-op", base, []string{"-a", "/nope"}, base},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.env)
			got, err := execCmd(t, newRemoveCmd(), tt.args...)
			if err != nil {
				t.Fatalf("remove %v: unexpected error: %v", tt.args, err)
			}
			if got != tt.want {
				t.Errorf("remove %v = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestRemoveCmd_strictAbsentErrors(t *testing.T) {
	t.Setenv("PATH", strings.Join([]string{"/usr/bin", "/bin"}, pathSep))
	if _, err := execCmd(t, newRemoveCmd(), "/nope"); err == nil {
		t.Fatal("remove /nope (strict default): expected error, got nil")
	}
}

func TestRemoveCmd_mutuallyExclusiveFlags(t *testing.T) {
	t.Setenv("PATH", strings.Join([]string{"/a", "/x"}, pathSep))
	if _, err := execCmd(t, newRemoveCmd(), "-i", "-a", "/a"); err == nil {
		t.Fatal("remove -i -a: expected error, got nil")
	}
}

func TestRemoveCmd_requiresExactlyOneArg(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")
	if _, err := execCmd(t, newRemoveCmd()); err == nil {
		t.Fatal("remove with no args: expected error, got nil")
	}
}
