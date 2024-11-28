package functions

import "context"

type SubmitFilesServiceInput struct {
	BaseBranch string
}

type SubmitFilesCallerType func(input SubmitFilesInput) error

type SubmitFilesService interface {
	Caller(ctx context.Context, input SubmitFilesServiceInput) SubmitFilesCallerType
}
