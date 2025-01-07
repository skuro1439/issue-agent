package cli

import (
	"fmt"
	"os"

	"github/clover0/github-issue-agent/logger"
)

func Parse() (string, []string, error) {
	if len(os.Args) < 2 {
		return "", nil, fmt.Errorf("command is required")
	}

	return os.Args[1], os.Args[2:], nil
}

func Execute(lo logger.Logger) error {
	command, flags, err := Parse()
	if err != nil {
		lo.Error("failed to parse input: %s\n", err)
		os.Exit(1)
	}

	// TODO: bind common flags to common struct here

	switch command {
	case "issue":
		return IssueCommand(lo, flags)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

}
