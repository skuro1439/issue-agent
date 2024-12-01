package models

import (
	"context"
	"os"

	"github/clover0/github-issue-agent/agent"
	"github/clover0/github-issue-agent/functions"
	"github/clover0/github-issue-agent/logger"
	"github/clover0/github-issue-agent/step"
)

type OpenAILLMForwarder struct {
	openai OpenAI
}

func NewOpenAILLMForwarder(l logger.Logger) agent.LLMForwarder {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		panic("OPENAI_API_KEY is not set")
	}

	return OpenAILLMForwarder{
		openai: NewOpenAI(l, apiKey),
	}
}

func (o OpenAILLMForwarder) StartForward(input agent.StartCompletionInput) ([]agent.LLMMessage, error) {
	return o.openai.StartCompletion(
		context.TODO(),
		agent.StartCompletionInput{
			Model:           input.Model,
			SystemPrompt:    input.SystemPrompt,
			StartUserPrompt: input.StartUserPrompt,
			Functions:       functions.AllFunctions(),
		},
	)
}

func (o OpenAILLMForwarder) ForwardLLM(
	ctx context.Context,
	input agent.StartCompletionInput,
	llmContexts []step.ReturnToLLMContext,
	history []agent.LLMMessage,
) ([]agent.LLMMessage, error) {
	return o.openai.ContinueCompletion(ctx, input, llmContexts, history)
}

func (o OpenAILLMForwarder) ForwardStep(ctx context.Context, history []agent.LLMMessage) step.Step {
	return o.openai.CompletionNextStep(ctx, history)
}
