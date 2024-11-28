package loader

import (
	"context"

	"github.com/google/go-github/v66/github"
)

type GitHub struct {
	client *github.Client
}

func NewGitHub(gh *github.Client) Loader {
	return GitHub{gh}
}

func (g GitHub) GetIssue(ctx context.Context, owner, repo string, number int) (Issue, error) {
	issue, _, err := g.client.Issues.Get(ctx, owner, repo, number)
	if err != nil {
		return Issue{}, err
	}

	return Issue{
		Content: *issue.Body,
	}, nil
}
