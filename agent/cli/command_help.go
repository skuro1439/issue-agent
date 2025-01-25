package cli

import (
	"flag"
	"fmt"

	"github.com/clover0/issue-agent/logger"
)

func Help(lo logger.Logger) {
	msg := `Usage
  issue-agent <command> [flags]
Commands  
  Help: Show usage of commands and flags
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
