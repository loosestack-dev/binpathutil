package subcmd

import (
	"github.com/spf13/cobra"
)

// Version is the binary's version. It defaults to "dev" and is overridden at
// release time via -ldflags "-X .../subcmd.Version=<tag>" (see .goreleaser.yaml).
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "binpath",
	Short: "Manipulate and query your PATH env variable",
}

func Execute() error {
	rootCmd.Version = Version
	return rootCmd.Execute()
}
