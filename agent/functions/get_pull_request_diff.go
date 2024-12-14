package functions

import (
	"reflect"
)

const FuncGetPullRequestDiff = "get_pull_request_diff"

type RepositoryService interface {
	GetPullRequestDiff(prNumber string) (string, error)
}

type GetPullRequestDiffType func(input GetPullRequestDiffInput) (string, error)

func InitGetPullRequestFunction(service RepositoryService) Function {
	f := Function{
		Name:        FuncGetPullRequestDiff,
		Description: "Get a Pull Request diff",
		Func:        GetPullRequestDiffCaller(service),
		FuncType:    reflect.TypeOf(GetPullRequestDiffCaller(service)),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pr_number": map[string]interface{}{
					"type":        "string",
					"description": "To get a Pull Request Number",
				},
			},
			"required":             []string{"pr_number"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncGetPullRequestDiff] = f

	return f
}

type GetPullRequestDiffInput struct {
	PRNumber string `json:"pr_number"`
}

func GetPullRequestDiffCaller(service RepositoryService) GetPullRequestDiffType {
	return func(input GetPullRequestDiffInput) (string, error) {
		return service.GetPullRequestDiff(input.PRNumber)
	}
}
