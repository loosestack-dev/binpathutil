package cmd 

import (
	"fmt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use: "add",
	Short: "Add an entry in the PATH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("add")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
