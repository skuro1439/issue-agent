package models

import (
	"fmt"
	"strings"

	"github/clover0/github-issue-agent/logger"
)

func SelectForwarder(lo logger.Logger, model string) LLMForwarder {
	if strings.HasPrefix(model, "gpt") {
		return NewOpenAILLMForwarder(lo)
	}

	if strings.HasPrefix(model, "claude") {
		return NewAnthropicLLMForwarder(lo)
	}

	panic(fmt.Sprintf("model %s is not supported\n", model))
}
