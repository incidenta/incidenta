package main

import (
	"os"

	"github.com/incidenta/incidenta/cmd/incidenta/command"
)

func main() {
	if err := command.NewRootCmd(os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
