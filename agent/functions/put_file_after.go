package functions

import (
	"github/clover0/github-issue-agent/store"
)

// TODO: implement hook?
func StoreFileAfterPutFile(s *store.Store, file store.File) {
	s.AddChangedFile(file)
}
