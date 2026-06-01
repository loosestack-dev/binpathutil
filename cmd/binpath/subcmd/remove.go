package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"binpathutil/internal/path"
	"binpathutil/internal/path/usecases"
)

func newRemoveCmd() *cobra.Command {
	var (
		ifPresent bool
		all       bool
	)

	cmd := &cobra.Command{
		Use:   "remove <entry>",
		Short: "Remove an entry from the PATH",
		Long: "Remove an entry from the PATH and print the resulting value.\n\n" +
			"By default removing an entry that is not present is an error. Pass\n" +
			"--if-present to make it a no-op, or --all to remove every occurrence.\n\n" +
			"The current process cannot modify the parent shell's PATH, so the new\n" +
			"value is written to stdout. Apply it with: export PATH=$(binpath remove <entry>)",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			element := args[0]

			var remove func(string, func() (string, error)) (string, error)
			switch {
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
	return cmd
}

func init() {
	rootCmd.AddCommand(newRemoveCmd())
}
