package subcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/loosestack-dev/binpathutil/internal/path"
	"github.com/loosestack-dev/binpathutil/internal/path/usecases"
)

var ErrNotPresent = errors.New("entry not present in PATH")

func newContainsCmd() *cobra.Command {
	var useRegex bool

	cmd := &cobra.Command{
		Use:   "contains <entry|pattern>",
		Short: "Report, via the exit code, whether the PATH contains an entry",
		Long: "Check whether the PATH contains the given entry.\n\n" +
			"By default the argument is matched literally; pass --regex to treat it\n" +
			"as a regular expression matched against each PATH entry.\n\n" +
			"Produces no output: the result is communicated through the exit code,\n" +
			"so it composes in shell conditionals like test(1) or grep -q:\n\n" +
			"  if binpath contains /usr/bin; then echo present; fi\n\n" +
			"Exit code is 0 when present, 1 when absent, 2 on error (e.g. a bad pattern).",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			argument := args[0]

			// Both use cases share the same signature, so --regex just selects which.
			check := usecases.Contains
			if useRegex {
				check = usecases.ContainsRegex
			}

			found, err := check(argument, path.GetEnvPath)
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

	cmd.Flags().BoolVarP(&useRegex, "regex", "r", false,
		"treat the argument as a regular expression matched against each PATH entry")
	return cmd
}

func init() {
	rootCmd.AddCommand(newContainsCmd())
}
