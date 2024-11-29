package functions

import (
	"fmt"
	"os"
	"path/filepath"
)

const FuncPutFile = "put_file"

func NewPutFileFunction() Function {
	f := Function{
		Name:        FuncPutFile,
		Description: "Put new content to the file",
		Func:        PutFile,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"output_path": map[string]interface{}{
					"type":        "string",
					"description": "Path of the file to be changed to the new content",
				},
				"content_text": map[string]interface{}{
					"type":        "string",
					"description": "The new content of the file",
				},
			},
			"required":             []string{"output_path", "content_text"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncPutFile] = f

	return f
}

type PutFileInput struct {
	OutputPath  string `json:"output_path"`
	ContentText string `json:"content_text"`
}

func PutFile(input PutFileInput) error {
	baseDir := filepath.Dir(input.OutputPath)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("mkdir all %s error: %w", baseDir, err)
	}

	f, err := os.Create(input.OutputPath)
	if err != nil {
		return fmt.Errorf("putting %s: %w", input.OutputPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(input.ContentText); err != nil {
		return err
	}

	return nil
}
