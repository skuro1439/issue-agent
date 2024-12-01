package agent

import (
	"context"
	"fmt"
	"github/clover0/github-issue-agent/store"

	"github/clover0/github-issue-agent/functions"
	"github/clover0/github-issue-agent/logger"
	"github/clover0/github-issue-agent/prompt"
	"github/clover0/github-issue-agent/step"
)

type AgentLike interface {
	Work()
}

type Agent struct {
	name                string
	parameter           Parameter
	currentStep         step.Step
	logg                logger.Logger
	submitServiceCaller functions.SubmitFilesCallerType
	llmForwarder        LLMForwarder
	prompt              prompt.Prompt
	history             []LLMMessage
	store               *store.Store
}

func NewAgent(
	parameter Parameter,
	name string,
	logg logger.Logger,
	submitServiceCaller functions.SubmitFilesCallerType,
	prompt prompt.Prompt,
	forwarder LLMForwarder,
	store *store.Store,
) Agent {
	return Agent{
		name:                name,
		parameter:           parameter,
		currentStep:         step.Step{},
		logg:                logg,
		submitServiceCaller: submitServiceCaller,
		prompt:              prompt,
		llmForwarder:        forwarder,
		store:               store,
	}
}

func (a *Agent) Work() (lastOutput string, err error) {
	ctx := context.Background()
	a.logg.Info("[%s]start agent work\n", a.name)

	completionInput := StartCompletionInput{
		Model:           a.parameter.Model,
		SystemPrompt:    a.prompt.SystemPrompt,
		StartUserPrompt: a.prompt.StartUserPrompt,
		Functions:       functions.AllFunctions(),
	}

	history, err := a.llmForwarder.StartForward(completionInput)
	if err != nil {
		return lastOutput, fmt.Errorf("start llm forward error: %w", err)
	}
	a.updateHistory(history)

	a.currentStep = a.llmForwarder.ForwardStep(ctx, history)

	var i int
	loop := true
	for loop {
		i++
		if i > a.parameter.MaxSteps {
			a.logg.Info("Reached to the max steps")
			break
		}

		switch a.currentStep.Do {
		case step.Exec:
			var input []step.ReturnToLLMInput
			for _, fnCtx := range a.currentStep.FunctionContexts {
				var returningStr string
				returningStr, err = functions.ExecFunction(
					a.store,
					fnCtx.Function.Name,
					fnCtx.FunctionArgs.String(),
					functions.SetSubmitFiles(
						a.submitServiceCaller,
					),
				)

				if err != nil {
					returningStr = "error caused! error message is: " + err.Error()
				}

				input = append(input, step.ReturnToLLMInput{
					ToolCallerID: fnCtx.ToolCallerID,
					Content:      returningStr,
				})
			}
			a.currentStep = step.NewReturnToLLMStep(input)

		case step.ReturnToLLM:
			history, err = a.llmForwarder.ForwardLLM(ctx, completionInput, a.currentStep.ReturnToLLMContexts, history)
			if err != nil {
				a.logg.Error("unrecoverable ContinueCompletion: %s\n", err)
				return lastOutput, err
			}
			a.updateHistory(history)
			a.currentStep = a.llmForwarder.ForwardStep(ctx, history)
			a.logg.Debug("end step return to LLM\n")

		case step.WaitingInstruction:
			a.logg.Debug("finish instruction\n")
			lastOutput = a.currentStep.LastOutput
			loop = false
			break

		case step.Unrecoverable, step.Unknown:
			a.logg.Error("unrecoverable error: %s\n", a.currentStep.UnrecoverableErr)
			return lastOutput, err
		default:
			a.logg.Error("does not exist step type\n")
			return lastOutput, err
		}
	}
	return lastOutput, nil
}

func (a *Agent) updateHistory(history []LLMMessage) {
	a.history = history
}

func (a *Agent) History() []LLMMessage {
	return a.history
}

func (a *Agent) ChangedFiles() []store.File {
	return a.store.ChangedFiles()
}
