package functions

import "github/clover0/github-issue-agent/store"

func StoreFileAfterModifyFile(s *store.Store, file store.File) {
	StoreFileAfterPutFile(s, file)
}
