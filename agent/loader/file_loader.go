package loader

import (
	"context"
	"io"
	"log"
	"os"
)

type FileLoader struct{}

func NewFileLoader() Loader {
	return &FileLoader{}
}

func (f *FileLoader) LoadIssue(_ context.Context, path string) (Issue, error) {
	file, err := os.Open(path)
	if err != nil {
		return Issue{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return Issue{
		Path:    path,
		Content: string(data),
	}, nil
}
