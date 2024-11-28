package loader

import "context"

type Loader interface {
	GetIssue(ctx context.Context, owner, repo string, number int) (Issue, error)
}
