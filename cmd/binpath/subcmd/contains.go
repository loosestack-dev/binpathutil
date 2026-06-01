package subcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"binpathutil/internal/path"
	"binpathutil/internal/path/usecases"
)

var ErrNotPresent = errors.New("entry not present in PATH")

func newContainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contains <entry>",
		Short: "Report, via the exit code, whether the PATH contains an entry",
		Long: "Check whether the PATH contains the given entry.\n\n" +
			"Produces no output: the result is communicated through the exit code,\n" +
			"so it composes in shell conditionals like test(1) or grep -q:\n\n" +
			"  if binpath contains /usr/bin; then echo present; fi\n\n" +
			"Exit code is 0 when the entry is present, non-zero when it is absent.",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			element := args[0]

			found, err := usecases.Contains(element, path.GetEnvPath)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)
				return err
			}
			if !found {
				return ErrNotPresent
			}
			return nil
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newContainsCmd())
}
