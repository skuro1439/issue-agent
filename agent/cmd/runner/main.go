package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/google/go-github/v66/github"

	"github/clover0/issue-agent/functions"
	"github/clover0/issue-agent/functions/agithub"
	"github/clover0/issue-agent/loader"
	"github/clover0/issue-agent/logger"
	"github/clover0/issue-agent/models"
	"github/clover0/issue-agent/prompt"
	"github/clover0/issue-agent/step"
)

func newOpenAI(l logger.Logger) models.OpenAI {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		panic("OPENAI_API_KEY is not set")
	}

	return models.NewOpenAI(l, apiKey)
}

func newGitHub() *github.Client {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		panic("GITHUB_TOKEN is not set")
	}
	return github.NewClient(nil).WithAuthToken(token)
}

func main() {
	workdir := os.Getenv("AGENT_WORKDIR")
	if err := os.Chdir(workdir); err != nil {
		log.Fatalf("failed to change directory: %s", err)
	}

	owner, repository := "clover0", "hobby"

	//lo := logger.NewDefaultLogger()
	lo := logger.NewPrinter()

	var temp string

	flag.StringVar(&temp, "template", "", "prompt template path")
	flag.Parse()

	y, err := prompt.LoadPromptTemplateFromYAML(temp)
	if err != nil {
		panic(err)
	}

	gh := newGitHub()

	issLoader := loader.NewGitHub(gh)
	ctx := context.Background()
	iss, err := issLoader.GetIssue(ctx, owner, repository, 5)
	if err != nil {
		panic(err)
	}

	m := map[string]string{
		"issue": iss.Content,
	}
	tpl, err := template.New("prompt").Parse(y.UserTemplate)
	if err != nil {
		panic(err)
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(tplbuff, m); err != nil {
		panic(err)
	}

	oai := newOpenAI(lo)
	chat, params, err := oai.StartCompletion(ctx, string(tplbuff.Bytes()), functions.AllFunctions())
	if err != nil {
		panic(err)
	}
	nextStep := oai.CompletionNextStep(ctx, chat)

	githubService := agithub.NewSubmitFileGitHubService(
		owner, repository, gh, lo,
	)

	var i int
	for {
		i++
		if i > 10 {
			fmt.Println("Reached to the max steps")
			break
		}
		fmt.Println("next step")

		switch nextStep.Do {
		case step.Exec:
			var input []step.ReturnToLLMInput
			for _, fnCtx := range nextStep.FunctionContexts {
				str, err := functions.ExecFunction(
					fnCtx.Function.Name,
					fnCtx.FunctionArgs.String(),
					functions.SetSubmitFiles(
						githubService.Caller(ctx,
							functions.SubmitFilesServiceInput{
								BaseBranch: "main", // TODO: changeable
							},
						),
					),
				)

				if err != nil {
					log.Fatalf("unrecoverable ExecFunction: %s", err)
				}
				input = append(input, step.ReturnToLLMInput{
					ToolCallerID: fnCtx.ToolCallerID,
					Content:      str,
				})
			}
			nextStep = step.NewReturnToLLMStep(input)

			log.Println("end step exec")

		case step.ReturnToLLM:
			chat, err = oai.ContinueCompletion(ctx, *chat, nextStep.ReturnToLLMContexts, &params)
			if err != nil {
				log.Fatalf("unrecoverable ContinueCompletion: %s", err)
			}
			nextStep = oai.CompletionNextStep(ctx, chat)
			log.Println("end step return to LLM")

		case step.WaitingInstruction:
			fmt.Println("finish instruction")
			break
		case step.Unrecoverable, step.Unknown:
			log.Fatalf("unrecoverable step type")
		default:
			log.Fatalf("does not exist step type")
		}

		fmt.Println("end step")
	}

	fmt.Println("Agent finished successfully!")

}
