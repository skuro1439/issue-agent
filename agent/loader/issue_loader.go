package loader

import "context"

type Loader interface {
	LoadIssue(ctx context.Context, number string) (Issue, error)
}
