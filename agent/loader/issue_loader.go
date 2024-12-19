package loader

import "context"

type Loader interface {
	LoadIssue(ctx context.Context, path string) (Issue, error)
}
