package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"binpathutil/internal/path"
	"binpathutil/internal/path/usecases"
)

func newAddCmd() *cobra.Command {
	var (
		ifAbsent bool
		first    bool
		last     bool
	)

	cmd := &cobra.Command{
		Use:   "add <entry>",
		Short: "Add an entry to the PATH",
		Long: "Add an entry to the PATH and print the resulting value.\n\n" +
			"By default the entry is added at the front (--first); pass --last to\n" +
			"append it instead.\n\n" +
			"The current process cannot modify the parent shell's PATH, so the new\n" +
			"value is written to stdout. Apply it with: export PATH=$(binpath add <entry>)",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			element := args[0]

			var add func(string, func() (string, error)) (string, error)
			switch {
			case last && ifAbsent:
				add = usecases.AppendIfAbsent
			case last:
				add = usecases.Append
			case ifAbsent:
				add = usecases.PrependIfAbsent
			default:
				add = usecases.Prepend
			}

			newPath, err := add(element, path.GetEnvPath)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), newPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&first, "first", "f", false, "add the entry at the front of PATH (default)")
	cmd.Flags().BoolVarP(&last, "last", "l", false, "add the entry at the end of PATH")
	cmd.MarkFlagsMutuallyExclusive("first", "last")
	cmd.Flags().BoolVarP(&ifAbsent, "if-absent", "i", false,
		"only add the entry if it is not already present in PATH")
	return cmd
}

func init() {
	rootCmd.AddCommand(newAddCmd())
}
