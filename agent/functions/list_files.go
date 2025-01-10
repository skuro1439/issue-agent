package functions

import (
	"fmt"
	"os"
)

const FuncListFiles = "list_files"

func InitListFilesFunction() Function {
	f := Function{
		Name:        FuncListFiles,
		Description: "List the files within the directory like Unix ls command. Each line contains the file mode, byte size, and name",
		Func:        ListFiles,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The valid path to list within its directory",
				},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncListFiles] = f

	return f
}

type ListFilesInput struct {
	Path string
}

func ListFiles(input ListFilesInput) ([]string, error) {
	if err := guardPath(input.Path); err != nil {
		return nil, err
	}

	if _, err := os.Stat(input.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist: %w", input.Path, err)
	}

	entries, err := os.ReadDir(input.Path)
	if err != nil {
		return nil, fmt.Errorf("can't read directory at %sr: %w", input.Path, err)
	}

	var files []string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info error: %w", err)
		}

		files = append(files,
			fmt.Sprintf("%s %d %s", info.Mode(), info.Size(), entry.Name()),
		)
	}

	return files, nil
}
