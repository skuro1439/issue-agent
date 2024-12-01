package store

type Store struct {
	changedFiles []File
}

func NewStore() Store {
	return Store{
		changedFiles: []File{},
	}
}

func (s *Store) AddChangedFile(f File) {
	s.changedFiles = append(s.changedFiles, f)
}

func (s *Store) ChangedFiles() []File {
	return s.changedFiles
}
