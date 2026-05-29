package main

import (
	"os"

	"binpathutil/cmd/binpath/subcmd"
)

func main() {
	if err := subcmd.Execute(); err != nil {
		os.Exit(1)
	}
}
