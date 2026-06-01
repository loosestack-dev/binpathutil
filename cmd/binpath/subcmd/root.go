package subcmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "binpath",
	Short: "Manipulate and query your PATH env variable",
}

func Execute() error {
	return rootCmd.Execute()
}
