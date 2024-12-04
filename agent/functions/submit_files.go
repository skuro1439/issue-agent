package functions

const FuncSubmitFiles = "submit_files"

func InitSubmitFilesGitHubFunction() Function {
	// TODO: selectable other method

	f := Function{
		Name:        FuncSubmitFiles,
		Description: "Submit the modified files by GitHub Pull Request",
		Func:        SubmitFiles,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"commit_message": map[string]interface{}{
					"type":        "string",
					"description": "Commit message indicating changes to the file",
				},
				"pull_request_content": map[string]interface{}{
					"type":        "string",
					"description": "Pull Request Content",
				},
			},
			"required":             []string{"commit_message", "pull_request_content"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncSubmitFiles] = f

	return f
}

type SubmitFilesInput struct {
	CommitMessage      string `json:"commit_message"`
	PullRequestContent string `json:"pull_request_content"`
}

func SubmitFiles(submitting func(input SubmitFilesInput) error, input SubmitFilesInput) error {
	return submitting(input)
}
