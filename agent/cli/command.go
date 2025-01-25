package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/clover0/issue-agent/logger"
)

// Parse parses the command and others from os.Args
// issue-agent <command> others
func Parse() (command string, others []string, err error) {
	if len(os.Args) < 2 {
		return "", nil, fmt.Errorf("command is required")
	}

	return os.Args[1], os.Args[2:], nil
}

func Execute() error {
	command, others, err := Parse()
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	// TODO: bind common flags to common struct here

	lo := logger.NewPrinter("info")
	switch command {
	case "version":
		return VersionCommand()
	case CreatePrCommand:
		return CreatePRCommand(others)
	case "help":
		help(lo)
		return nil
	default:
		help(lo)
		return fmt.Errorf("unknown command: %s", command)
	}

}

func help(lo logger.Logger) {
	msg := `Usage
  issue-agent <command> [flags]
Commands  help: Show usage of commands and flags
  help: Show usage of commands and flags
  version: Show version of issue-agent CLI
`
	issueFlags, _ := CreatePRFlags()
	msg += fmt.Sprintf("  %s:\n", CreatePrCommand)
	issueFlags.VisitAll(func(flg *flag.Flag) {
		msg += fmt.Sprintf("    --%s\n", flg.Name)
		msg += fmt.Sprintf("        %s\n", flg.Usage)
	})
	lo.Info(msg)
}
