package cli

import (
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

	lo := logger.NewPrinter("info")
	switch command {
	case "version":
		return VersionCommand()
	case CreatePrCommand:
		return CreatePRCommand(others)
	case "Help":
		Help(lo)
		return nil
	default:
		Help(lo)
		return fmt.Errorf("unknown command: %s", command)
	}
}
