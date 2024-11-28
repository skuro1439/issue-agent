package functions

import (
	"fmt"
	"os"
)

const FuncModifyFile = "modify_file"

func NewModifyFileFunction() Function {
	f := Function{
		Name:        FuncModifyFile,
		Description: "Modify the file at output_path with the contents of content_text",
		Func:        ModifyFile,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"output_path": map[string]interface{}{
					"type":        "string",
					"description": "Path of the file to be modified to the new content",
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

	functionsMap[FuncModifyFile] = f

	return f
}

type ModifyFileInput struct {
	OutputPath  string `json:"output_path"`
	ContentText string `json:"content_text"`
}

func ModifyFile(input ModifyFileInput) error {
	f, err := os.Create(input.OutputPath)
	if err != nil {
		return fmt.Errorf("modify %s: %w", input.OutputPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(input.ContentText); err != nil {
		return err
	}

	return nil
}
