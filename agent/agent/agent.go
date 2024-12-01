package agent

import (
	"context"
	"fmt"
	"log"

	"github/clover0/issue-agent/functions"
	"github/clover0/issue-agent/logger"
	"github/clover0/issue-agent/prompt"
	"github/clover0/issue-agent/step"
)

type AgentLike interface {
	Work()
}

type Agent struct {
	parameter           Parameter
	currentStep         step.Step
	logg                logger.Logger
	submitServiceCaller functions.SubmitFilesCallerType
	llmForwarder        LLMForwarder
	goal                string
	prompt              prompt.Prompt
}

func NewAgent(
	parameter Parameter,
	logg logger.Logger,
	submitServiceCaller functions.SubmitFilesCallerType,
	prompt prompt.Prompt,
	forwarder LLMForwarder,
) Agent {
	return Agent{
		parameter:           parameter,
		currentStep:         step.Step{},
		logg:                logg,
		submitServiceCaller: submitServiceCaller,
		prompt:              prompt,
		llmForwarder:        forwarder,
	}
}

func (a *Agent) Work() error {
	ctx := context.Background()

	completionInput := StartCompletionInput{
		Model:           a.parameter.Model,
		SystemPrompt:    a.prompt.SystemPrompt,
		StartUserPrompt: a.prompt.StartUserPrompt,
		Functions:       functions.AllFunctions(),
	}

	history, err := a.llmForwarder.StartForward(completionInput)
	if err != nil {
		return fmt.Errorf("start llm forward error: %w", err)
	}

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
			// ForwardExec
			var input []step.ReturnToLLMInput
			for _, fnCtx := range a.currentStep.FunctionContexts {
				var returningStr string
				returningStr, err = functions.ExecFunction(
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
				log.Fatalf("unrecoverable ContinueCompletion: %s", err)
			}
			a.currentStep = a.llmForwarder.ForwardStep(ctx, history)
			a.logg.Debug("end step return to LLM")

		case step.WaitingInstruction:
			a.logg.Debug("finish instruction")
			loop = false
			break

		case step.Unrecoverable, step.Unknown:
			log.Fatalf(fmt.Sprintf("unnrecoverable error: %s", a.currentStep.UnrecoverableErr))
		default:
			log.Fatalf("does not exist step type")
		}
	}
	return nil
}
