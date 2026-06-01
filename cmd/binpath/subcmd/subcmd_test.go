package subcmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

const pathSep = string(os.PathListSeparator)

// Helper command for the different subcommand test, to pass the arguments
// as if they were given in a CLI
func execCmd(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return strings.TrimSpace(buf.String()), err
}
