package main

import (
	"os"
	"binpathutil/cmd/binpath/cmd"
) 

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
