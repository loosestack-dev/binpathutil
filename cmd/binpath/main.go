package main

import (
	"errors"
	"os"

	"github.com/loosestack-dev/binpathutil/cmd/binpath/subcmd"
)

// exit code mapping to have different result based on "the PATH does not countains what you asked for"
// and "there was an acutal error".
func exitCode(err error) int {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, subcmd.ErrNotPresent):
		return 1
	default:
		return 2
	}
}

func main() {
	os.Exit(exitCode(subcmd.Execute()))
}
