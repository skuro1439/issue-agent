package template

import _ "embed"

//go:embed default_prompt_ja.yaml
var defaultJAPrompt []byte

func DefaultTemplate() []byte {
	return defaultJAPrompt
}
