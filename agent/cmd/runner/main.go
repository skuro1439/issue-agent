package main

import (
	"context"
	"os"

	"github.com/google/go-github/v66/github"

	"github/clover0/github-issue-agent/agent"
	"github/clover0/github-issue-agent/config/cli"
	"github/clover0/github-issue-agent/functions"
	"github/clover0/github-issue-agent/functions/agithub"
	"github/clover0/github-issue-agent/loader"
	"github/clover0/github-issue-agent/logger"
	"github/clover0/github-issue-agent/models"
	libprompt "github/clover0/github-issue-agent/prompt"
	"github/clover0/github-issue-agent/store"
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

	llmForwarder := models.NewAnthropicLLMForwarder(lo)

	cliIn, err := cli.ParseInput()
	if err != nil {
		lo.Error("failed to parse input: %s", err)
		os.Exit(1)
	}

	if err := agithub.CloneRepository(lo, cliIn); err != nil {
		lo.Error("failed to clone repository")
		os.Exit(1)
	}

	// TODO: no dependency with changing directory
	if err := os.Chdir(cliIn.AgentWorkDir); err != nil {
		lo.Error("failed to change directory: %s", err)
		os.Exit(1)
	}

	promptTemplate, err := libprompt.LoadPromptTemplateFromYAML(cliIn.Template)
	if err != nil {
		lo.Error("failed to load prompt template: %s", err)
		os.Exit(1)
	}

	ctx := context.Background()

	gh := newGitHub()
	issLoader := loader.NewGitHubLoader(gh, cliIn.RepositoryOwner, cliIn.Repository)
	submitServiceCaller := agithub.NewSubmitFileGitHubService(cliIn.RepositoryOwner, cliIn.Repository, gh, lo).
		Caller(ctx, functions.SubmitFilesServiceInput{
			BaseBranch: cliIn.BaseBranch,
			GitEmail:   cliIn.GitEmail,
			GitName:    cliIn.GitName,
		})

	dataStore := store.NewStore()

	parameter := agent.Parameter{
		MaxSteps: cliIn.MaxSteps,
		Model:    cliIn.Model,
	}
	requirementAgent := RunRequirementAgent(promptTemplate, issLoader, submitServiceCaller, parameter, cliIn, lo, &dataStore, llmForwarder)

	//developerAgent := RunDeveloperAgent(promptTemplate, issLoader, cliIn, lo, gh, &dataStore, llmForwarder)
	developer2Agent := RunDeveloper2Agent(promptTemplate, issLoader, submitServiceCaller, parameter, cliIn, lo, &dataStore,
		requirementAgent.History()[len(requirementAgent.History())-1].RawContent, llmForwarder,
	)

	RunSecurityAgent(promptTemplate, developer2Agent.ChangedFiles(), submitServiceCaller, parameter, lo, &dataStore, llmForwarder)

	lo.Info("Agents finished successfully!")
}

func RunRequirementAgent(
	promptTemplate libprompt.PromptTemplate,
	issLoader loader.Loader,
	submitServiceCaller functions.SubmitFilesCallerType,
	parameter agent.Parameter,
	cliIn cli.Inputs,
	lo logger.Logger,
	dataStore *store.Store,
	llmForwarder agent.LLMForwarder,
) agent.Agent {
	prompt, err := libprompt.BuildRequirementPrompt(promptTemplate, issLoader, cliIn.GithubIssueNumber)
	if err != nil {
		lo.Error("failed buld prompt: %s", err)
		os.Exit(1)
	}

	ag := agent.NewAgent(
		parameter,
		"main",
		lo,
		submitServiceCaller,
		prompt,
		llmForwarder,
		dataStore,
	)

	_, err = ag.Work()
	if err != nil {
		lo.Error("ag failed: %s", err)
		os.Exit(1)
	}

	return ag
}

func RunDeveloperAgent(
	promptTemplate libprompt.PromptTemplate,
	issLoader loader.Loader,
	submitServiceCaller functions.SubmitFilesCallerType,
	parameter agent.Parameter,
	cliIn cli.Inputs,
	lo logger.Logger,
	dataStore *store.Store,
	llmForwarder agent.LLMForwarder,
) agent.Agent {
	prompt, err := libprompt.BuildDeveloperPrompt(promptTemplate, issLoader, cliIn.GithubIssueNumber)
	if err != nil {
		lo.Error("failed build prompt: %s", err)
		os.Exit(1)
	}

	ag := agent.NewAgent(
		parameter,
		"main",
		lo,
		submitServiceCaller,
		prompt,
		llmForwarder,
		dataStore,
	)

	_, err = ag.Work()
	if err != nil {
		lo.Error("ag failed: %s", err)
		os.Exit(1)
	}

	return ag
}

func RunDeveloper2Agent(
	promptTemplate libprompt.PromptTemplate,
	issLoader loader.Loader,
	submitServiceCaller functions.SubmitFilesCallerType,
	parameter agent.Parameter,
	cliIn cli.Inputs,
	lo logger.Logger,
	dataStore *store.Store,
	instruction string,
	llmForwarder agent.LLMForwarder,
) agent.Agent {
	prompt, err := libprompt.BuildDeveloper2Prompt(promptTemplate, issLoader, cliIn.GithubIssueNumber, instruction)
	if err != nil {
		lo.Error("failed build prompt: %s", err)
		os.Exit(1)
	}

	ag := agent.NewAgent(
		parameter,
		"main",
		lo,
		submitServiceCaller,
		prompt,
		llmForwarder,
		dataStore,
	)

	_, err = ag.Work()
	if err != nil {
		lo.Error("developer agent failed: %s", err)
		os.Exit(1)
	}

	return ag
}

func RunSecurityAgent(
	promptTemplate libprompt.PromptTemplate,
	changedFiles []store.File,
	submitServiceCaller functions.SubmitFilesCallerType,
	parameter agent.Parameter,
	lo logger.Logger,
	dataStore *store.Store,
	llmForwarder agent.LLMForwarder,
) agent.Agent {
	var changedFilePath []string
	for _, f := range changedFiles {
		changedFilePath = append(changedFilePath, f.Path)
	}

	securityPrompt, err := libprompt.BuildSecurityPrompt(promptTemplate, changedFilePath)
	if err != nil {
		lo.Error("failed to build security prompt: %s", err)
		os.Exit(1)
	}
	ag := agent.NewAgent(
		parameter,
		"securityAgent",
		lo,
		submitServiceCaller,
		securityPrompt,
		llmForwarder,
		dataStore,
	)

	if _, err := ag.Work(); err != nil {
		lo.Error("securityAgent failed: %s", err)
		os.Exit(1)
	}

	return ag
}
