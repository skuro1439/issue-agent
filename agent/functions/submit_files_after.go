package functions

import "github/clover0/github-issue-agent/store"

func SubmitFilesAfter(s *store.Store, storeKey string, storeValue SubmitFilesOutput) {
	s.AddSubmission(storeKey, store.Submission{
		BaseBranch:        storeValue.Branch,
		PullRequestNumber: storeValue.PullRequestNumber,
	})
}
