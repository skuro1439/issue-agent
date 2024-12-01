package cli

import (
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Inputs struct {
	Model             string
	MaxSteps          int
	Template          string
	Repository        string `validate:"required"`
	RepositoryOwner   string `validate:"required"`
	GithubIssueNumber string `validate:"required"`
	BaseBranch        string `validate:"required"`
	AgentWorkDir      string
}

func ParseInput() (Inputs, error) {
	cliIn := Inputs{}

	flag.StringVar(&cliIn.Model, "model", "gpt-4o", "Prompt template path")
	flag.IntVar(&cliIn.MaxSteps, "max_steps", 100, "Max steps for the agent to run. Avoid infinite loop.")
	flag.StringVar(&cliIn.Template, "template", "./default_prompt_ja.yaml", "Prompt template path")
	flag.StringVar(&cliIn.RepositoryOwner, "repository_owner", "", "GitHubLoader Repository owner")
	flag.StringVar(&cliIn.Repository, "repository", "", "Working at GitHubLoader Repository name")
	flag.StringVar(&cliIn.GithubIssueNumber, "github_issue_number", "", "GitHubLoader issue number")
	flag.StringVar(&cliIn.BaseBranch, "base_branch", "", "Base Branch for pull request")
	flag.StringVar(&cliIn.AgentWorkDir, "workdir", "./", "Workdir for the agent to run")

	flag.Parse()

	validate := validator.New()
	if err := validate.Struct(cliIn); err != nil {
		errs := err.(validator.ValidationErrors)
		return Inputs{}, fmt.Errorf("validation failed: %w", errs)
	}

	return cliIn, nil
}
