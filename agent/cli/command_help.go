package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/clover0/issue-agent/logger"
)

const HelpCommand = "help"

func Help(lo logger.Logger) {
	msg := `Usage
  issue-agent <command> [flags]
Command and Flags  
  help: Show usage of commands and flags
  version: Show version of issue-agent CLI
`
	createPRFlags, _ := CreatePRFlags()
	msg += fmt.Sprintf("  %s:\n", CreatePrCommand)
	createPRFlags.VisitAll(func(flg *flag.Flag) {
		msg += fmt.Sprintf("    --%s\n", flg.Name)
		msg += IndentMultiLine(flg.Usage, "      ")
		msg += "\n"
	})
	lo.Info(msg)
}

func IndentMultiLine(str string, indent string) string {
	lines := strings.Split(str, "\n")
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = indent + line
	}

	return strings.Join(out, "\n")
}
