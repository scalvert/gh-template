package main

import (
	"os"

	"github.com/scalvert/gh-template/cmd"
)

func main() {
	templateCmd := cmd.NewCmdTemplate()

	if err := templateCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
