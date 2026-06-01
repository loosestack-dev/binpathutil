package subcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var containsCmd = &cobra.Command{
	Use: "contains",
	Short: "Verify if the the PATH contains the given element",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("contains")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(containsCmd)
}
