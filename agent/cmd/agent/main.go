package main

import (
	"os"

	"github/clover0/github-issue-agent/cli"
	"github/clover0/github-issue-agent/logger"
)

func main() {
	// TODO:
	//lo := logger.NewDefaultLogger()
	lo := logger.NewPrinter("error")

	if err := cli.Execute(lo); err != nil {
		lo.Error("failed to execute command: %s\n", err)
		os.Exit(1)
	}
}
