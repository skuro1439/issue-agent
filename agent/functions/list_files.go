package functions

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const FuncListFiles = "list_files"

func NewListFilesFunction() Function {
	f := Function{
		Name:        FuncListFiles,
		Description: "List the files within the directory",
		Func:        ListFiles,
		FuncType:    reflect.TypeOf(ListFiles),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"root": map[string]interface{}{
					"type":        "string",
					"description": "The path to list within its directory",
				},
			},
			"required":             []string{"root"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncListFiles] = f

	return f
}

type ListFilesInput struct {
	Root string
}

func ListFiles(input ListFilesInput) ([]string, error) {
	files := []string{}

	err := filepath.Walk(input.Root, func(path string, info os.FileInfo, err error) error {
		// TODO: selectable file type
		// ignore hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
