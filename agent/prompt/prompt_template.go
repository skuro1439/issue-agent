package prompt

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type PromptTemplate struct {
	SystemTemplate string `yaml:"system_template"`
	UserTemplate   string `yaml:"user_template"`
}

func LoadPromptTemplateFromYAML(filePath string) (PromptTemplate, error) {
	var pt PromptTemplate

	file, err := os.Open(filePath)
	if err != nil {
		return pt, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return pt, err
	}

	err = yaml.Unmarshal(data, &pt)
	if err != nil {
		return pt, err
	}

	return pt, nil
}
