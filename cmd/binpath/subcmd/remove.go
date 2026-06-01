package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an entry from the PATH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("remove")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
