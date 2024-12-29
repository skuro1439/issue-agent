package cli

import (
	"fmt"
	"os"
)

func Parse() ([]string, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("command is required")
	}

	return os.Args[2:], nil
}
