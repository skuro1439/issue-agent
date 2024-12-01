package main

import (
	"context"
	"github/clover0/issue-agent/config/cli"
	"github/clover0/issue-agent/models"
	"os"

	"github.com/google/go-github/v66/github"

	"github/clover0/issue-agent/agent"
	"github/clover0/issue-agent/functions"
	"github/clover0/issue-agent/functions/agithub"
	"github/clover0/issue-agent/loader"
	"github/clover0/issue-agent/logger"
	"github/clover0/issue-agent/prompt"
)

func newGitHub() *github.Client {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		panic("GITHUB_TOKEN is not set")
	}
	return github.NewClient(nil).WithAuthToken(token)
}

func main() {
	//lo := logger.NewDefaultLogger()
	lo := logger.NewPrinter()

	cliIn, err := cli.ParseInput()
	if err != nil {
		lo.Error("failed to parse input: %s", err)
		os.Exit(1)
	}

	if err := os.Chdir(cliIn.AgentWorkDir); err != nil {
		lo.Error("failed to change directory: %s", err)
		os.Exit(1)
	}

	promptTemplate, err := prompt.LoadPromptTemplateFromYAML(cliIn.Template)
	if err != nil {
		lo.Error("failed to load prompt template: %s", err)
		os.Exit(1)
	}

	gh := newGitHub()

	issLoader := loader.NewGitHubLoader(gh, cliIn.RepositoryOwner, cliIn.Repository)

	ctx := context.Background()

	prompt, err := prompt.BuildPrompt(promptTemplate, issLoader, cliIn.GithubIssueNumber)
	if err != nil {
		lo.Error("failed buld prompt: %s", err)
		os.Exit(1)
	}

	agent := agent.NewAgent(
		agent.Parameter{
			MaxSteps: cliIn.MaxSteps,
			Model:    cliIn.Model,
		},
		lo,
		agithub.NewSubmitFileGitHubService(cliIn.RepositoryOwner, cliIn.Repository, gh, lo).
			Caller(ctx, functions.SubmitFilesServiceInput{BaseBranch: cliIn.BaseBranch}),
		prompt,
		models.NewOpenAILLMForwarder(lo, prompt),
	)

	if err := agent.Work(); err != nil {
		lo.Error("agent failed: %s", err)
		os.Exit(1)
	}

	lo.Info("Agent finished successfully!")
}
