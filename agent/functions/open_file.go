package functions

import (
	"io"
	"log"
	"os"
)

const FuncOpenFile = "open_file"

func NewOpenFileFunction() Function {
	f := Function{
		Name:        FuncOpenFile,
		Description: "Open the file",
		Func:        OpenFile,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The path of the file to open",
				},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncOpenFile] = f

	return f
}

type File struct {
	Content string
}

type OpenFileInput struct {
	Path string
}

func OpenFile(input OpenFileInput) (File, error) {
	file, err := os.Open(input.Path)
	if err != nil {
		return File{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return File{Content: string(data)}, nil
}
