package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type IssueInputs struct {
	Model           string
	MaxSteps        int
	CloneRepository bool
	// TODO: make default value
	Template          string
	Repository        string `validate:"required"`
	RepositoryOwner   string `validate:"required"`
	GithubIssueNumber string
	BaseBranch        string `validate:"required"`
	GitName           string `validate:"required"`
	GitEmail          string `validate:"required"`
	AgentWorkDir      string
	FromFile          string
	NoSubmit          bool
}

// TODO: parse common inputs
func ParseIssueInput() (IssueInputs, error) {
	cliIn := IssueInputs{}

	cmd := flag.NewFlagSet("issue", flag.ExitOnError)
	// TODO: from yaml file
	cmd.StringVar(&cliIn.Model, "model", "gpt-4o", "Prompt template path")
	// TODO: from yaml file
	cmd.IntVar(&cliIn.MaxSteps, "max_steps", 100, "Max steps for the agent to run. Avoid infinite loop.")
	cmd.StringVar(&cliIn.Template, "template", "", "Prompt template path. default is `config/template/default_prompt_ja.yaml`")
	cmd.BoolVar(&cliIn.CloneRepository, "clone_repository", false, "Whether to clone repository to the workdir")
	cmd.StringVar(&cliIn.RepositoryOwner, "repository_owner", "", "GitHubLoader Repository owner")
	cmd.StringVar(&cliIn.Repository, "repository", "", "Working at GitHubLoader Repository name")
	cmd.StringVar(&cliIn.GithubIssueNumber, "github_issue_number", "", "GitHubLoader issue number")
	cmd.StringVar(&cliIn.BaseBranch, "base_branch", "", "Base Branch for pull request")
	cmd.StringVar(&cliIn.GitName, "git_name", "", "Name for git config using git commit")
	cmd.StringVar(&cliIn.GitEmail, "git_email", "", "Email for git config using git commit")
	// TODO: make default setting
	cmd.StringVar(&cliIn.AgentWorkDir, "workdir", ".", "Workdir for the agent to run")
	cmd.BoolVar(&cliIn.NoSubmit, "no_submit", false, "Option to not submit the pull request. And not commit.")
	cmd.StringVar(&cliIn.FromFile, "from_file", "", "Issue content from file path")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		return IssueInputs{}, fmt.Errorf("failed to parse input: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cliIn); err != nil {
		errs := err.(validator.ValidationErrors)
		return IssueInputs{}, fmt.Errorf("validation failed: %w", errs)
	}

	if cliIn.GithubIssueNumber == "" && cliIn.FromFile == "" {
		return IssueInputs{}, fmt.Errorf("github_issue_number or from_file is required")
	}

	return cliIn, nil
}
