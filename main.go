package main

import (
	"os"

	"github.com/forkspacer/cli/cmd"
	_ "github.com/forkspacer/cli/cmd/workspace"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
