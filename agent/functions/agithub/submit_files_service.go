package agithub

import (
	"context"
	"fmt"

	"github.com/google/go-github/v66/github"

	"github/clover0/issue-agent/functions"
	"github/clover0/issue-agent/functions/agit"
	"github/clover0/issue-agent/logger"
)

type SubmitFileGitHubService struct {
	owner      string
	repository string
	client     *github.Client
	logger     logger.Logger
}

func NewSubmitFileGitHubService(
	owner string,
	repository string,
	client *github.Client,
	logger logger.Logger,
) functions.SubmitFilesService {
	return SubmitFileGitHubService{
		owner:      owner,
		repository: repository,
		client:     client,
		logger:     logger,
	}
}

func (s SubmitFileGitHubService) Caller(
	ctx context.Context,
	callerInput functions.SubmitFilesServiceInput,
) functions.SubmitFilesCallerType {
	return func(input functions.SubmitFilesInput) error {

		var out string
		var err error

		out, err = agit.GitStatus()
		if err != nil {
			return fmt.Errorf("submit file service: git stastus error: %w", err)
		}

		newBranch := agit.MakeBranchName()

		out, err = agit.GitSwitchCreate(newBranch)
		if err != nil {
			return fmt.Errorf("submit file service: git switch error: %w", err)
		}
		s.logger.Debug(fmt.Sprintf("git swicth create: %s", out))

		out, err = agit.GitAddAll()
		if err != nil {
			return fmt.Errorf("submit file service: git add error: %w", err)
		}
		s.logger.Debug(fmt.Sprintf("git add all: %s\n", out))

		out, err = agit.GitCommit(input.CommitMessage)
		if err != nil {
			return fmt.Errorf("submit file service: git commit error: %w", err)
		}
		s.logger.Debug(fmt.Sprintf("git commit: %s\n", out))

		out, err = agit.GitPushBranch(newBranch)
		if err != nil {
			s.logger.Error(out)
			return fmt.Errorf("submit file service: git push branch error: %w", err)
		}

		s.logger.Debug(fmt.Sprintf("submit file service: create PR parameter: %s", callerInput))
		//head := fmt.Sprintf("%s:%s", *user.Name, newBranch)
		pr, _, err := s.client.PullRequests.Create(ctx, s.owner, s.repository, &github.NewPullRequest{
			Title: &input.CommitMessage,
			Head:  &newBranch,
			Base:  &callerInput.BaseBranch,
			Body:  &input.PullRequestContent,
		})
		if err != nil {
			return fmt.Errorf("submit file service: create PR: %w", err)
		}
		s.logger.Debug(fmt.Sprintf("created PR: %s", pr.URL))

		return nil
	}
}
