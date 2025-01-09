package cli

import (
	"fmt"
	"os"
)

func Parse() (command string, flags []string, err error) {
	if len(os.Args) < 2 {
		return "", nil, fmt.Errorf("command is required")
	}

	return os.Args[1], os.Args[2:], nil
}

func Execute() error {
	command, flags, err := Parse()
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	// TODO: bind common flags to common struct here

	switch command {
	case "issue":
		return IssueCommand(flags)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

}
