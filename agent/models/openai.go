package models

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github/clover0/issue-agent/functions"
	"github/clover0/issue-agent/logger"
	"github/clover0/issue-agent/step"
)

type OpenAI struct {
	client *openai.Client
	logger logger.Logger
}

func NewOpenAI(logger logger.Logger, apiKey string) OpenAI {
	return OpenAI{
		logger: logger,
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}

func (o OpenAI) StartCompletion(ctx context.Context, prompt string, functions []functions.Function) (*openai.ChatCompletion, openai.ChatCompletionNewParams, error) {
	toolFuncs := make([]openai.ChatCompletionToolParam, len(functions))
	for i, f := range functions {
		toolFuncs[i] = openai.ChatCompletionToolParam{
			Function: openai.F(f.ToFunctionCalling()),
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
		}

	}

	params := openai.ChatCompletionNewParams{
		// TODO: selectable model
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
		}),
		Temperature: openai.F(0.0),
		Tools:       openai.F(toolFuncs),
	}

	o.debugShowSendingMsg(params)
	chat, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, openai.ChatCompletionNewParams{}, err
	}

	o.logger.Debug(fmt.Sprintf("Prompt Token: %d, Completion Token: %d\n",
		chat.Usage.PromptTokens, chat.Usage.CompletionTokens,
	))
	o.logger.Debug(logger.Green("returned messages:\n"))
	o.debugShowChoice(chat.Choices)

	return chat, params, nil
}

func (o OpenAI) ContinueCompletion(ctx context.Context, completion openai.ChatCompletion, llmContexts []step.ReturnToLLMContext, params *openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	choice := completion.Choices[0]
	params.Messages.Value = append(params.Messages.Value, choice.Message)

	for _, v := range llmContexts {
		if v.ToolCallerID != "" {
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(v.ToolCallerID, v.Content))
		} else {
			params.Messages.Value = append(params.Messages.Value, openai.UserMessage(v.Content))
		}
	}

	o.debugShowSendingMsg(*params)
	chat, err := o.client.Chat.Completions.New(ctx, *params)
	if err != nil {
		return nil, err
	}

	o.logger.Debug(logger.Green("returned messages:\n"))
	o.debugShowChoice(chat.Choices)

	return chat, nil
}

func (o OpenAI) CompletionNextStep(_ context.Context, chat *openai.ChatCompletion) step.Step {
	choice := chat.Choices[0]
	if choice.FinishReason == openai.ChatCompletionChoicesFinishReasonToolCalls {
		var input []step.FunctionsInput
		for _, v := range choice.Message.ToolCalls {
			input = append(input, step.FunctionsInput{
				FuncName:     v.Function.Name,
				FunctionArgs: v.Function.Arguments,
				ToolCallerID: v.ID,
			})
		}
		return step.NewExecStep(input)
	}

	if choice.FinishReason == openai.ChatCompletionChoicesFinishReasonStop {
		return step.NewWaitingInstructionStep()
	}

	if choice.FinishReason == openai.ChatCompletionChoicesFinishReasonLength {
		return step.NewUnrecoverableStep()
	}

	return step.NewUnknownStep()
}

func (o OpenAI) debugShowSendingMsg(param openai.ChatCompletionNewParams) {
	if len(param.Messages.Value) > 0 {
		o.logger.Debug(logger.Green(fmt.Sprintf("model: %s, sending messages:\n", param.Model.String())))
		o.logger.Debug(fmt.Sprintf("%s\n", param.Messages.Value[len(param.Messages.Value)-1]))
	}
}

func (o OpenAI) debugShowChoice(completion []openai.ChatCompletionChoice) {
	for _, c := range completion {
		o.logger.Debug(fmt.Sprintf("finish_reason: %s, role: %s, message.content: %s\n",
			c.FinishReason, c.Message.Role, c.Message.Content,
		))
		o.logger.Debug("tools:\n")
		for _, t := range c.Message.ToolCalls {
			o.logger.Debug(fmt.Sprintf("id: %s, function_name:%s, function_args:%s\n", t.ID, t.Function.Name, t.Function.Arguments))
		}
	}
}
