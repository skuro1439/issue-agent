package functions

import (
	"fmt"
	"strings"
)

func guardPath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path %s contains not allowed '..'", path)
	}

	if strings.Contains(path, "~") {
		return fmt.Errorf("path %s contains not allowed '~'", path)
	}

	if strings.Contains(path, "//") {
		return fmt.Errorf("path %s contains not allowed '//'", path)
	}

	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("path %s starts with '/', not allowed", path)
	}

	return nil
}
