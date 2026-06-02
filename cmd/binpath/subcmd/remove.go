package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/loosestack-dev/binpathutil/internal/path"
	"github.com/loosestack-dev/binpathutil/internal/path/usecases"
)

func newRemoveCmd() *cobra.Command {
	var (
		ifPresent bool
		all       bool
		useRegex  bool
	)

	cmd := &cobra.Command{
		Use:   "remove <entry|pattern>",
		Short: "Remove an entry from the PATH",
		Long: "Remove an entry from the PATH and print the resulting value.\n\n" +
			"By default removing an entry that is not present is an error. Pass\n" +
			"--if-present to make it a no-op, or --all to remove every occurrence.\n" +
			"Pass --regex to match the argument as a regular expression against each\n" +
			"PATH entry instead of comparing it literally (composes with the above).\n\n" +
			"The current process cannot modify the parent shell's PATH, so the new\n" +
			"value is written to stdout. Apply it with: export PATH=$(binpath remove <entry>)",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			element := args[0]

			var remove func(string, func() (string, error)) (string, error)
			switch {
			case useRegex && all:
				remove = usecases.RemoveAllOccurenceRegex
			case useRegex && ifPresent:
				remove = usecases.RemoveIfPresentRegex
			case useRegex:
				remove = usecases.RemoveRegex
			case all:
				remove = usecases.RemoveAllOccurence
			case ifPresent:
				remove = usecases.RemoveIfPresent
			default:
				remove = usecases.Remove
			}

			newPath, err := remove(element, path.GetEnvPath)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), newPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&ifPresent, "if-present", "i", false,
		"do not error if the entry is not present in PATH")
	cmd.Flags().BoolVarP(&all, "all", "a", false,
		"remove every occurrence of the entry (strips duplicates)")
	cmd.MarkFlagsMutuallyExclusive("if-present", "all")
	cmd.Flags().BoolVarP(&useRegex, "regex", "r", false,
		"treat the argument as a regular expression matched against each PATH entry")
	return cmd
}

func init() {
	rootCmd.AddCommand(newRemoveCmd())
}
