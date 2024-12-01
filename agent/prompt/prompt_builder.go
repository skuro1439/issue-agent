package prompt

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github/clover0/github-issue-agent/loader"
)

type Prompt struct {
	SystemPrompt    string
	StartUserPrompt string
}

func BuildPrompt(promptTpl PromptTemplate, issueLoader loader.Loader, issueNumber string) (Prompt, error) {
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		panic(err)
	}

	m := map[string]string{
		"issue":       iss.Content,
		"issueNumber": issueNumber,
	}
	usrTpl, err := template.New("prompt").Parse(promptTpl.UserTemplate)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := usrTpl.Execute(tplbuff, m); err != nil {
		return Prompt{}, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return Prompt{
		SystemPrompt:    promptTpl.SystemTemplate,
		StartUserPrompt: tplbuff.String(),
	}, nil
}

// TODO: separeting prompt from yaml
func BuildSecurityPrompt(promptTpl PromptTemplate, changedFilesPath []string) (Prompt, error) {
	var prpt Prompt
	for _, p := range promptTpl.Agents {
		if p.Name == "security" {
			prpt = Prompt{
				SystemPrompt:    p.SystemTemplate,
				StartUserPrompt: p.UserTemplate,
			}
			break
		}
	}

	tpl, err := template.New("prompt").Parse(prpt.StartUserPrompt)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	m := map[string]any{
		"filePaths": changedFilesPath,
	}
	tplbuff := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(tplbuff, m); err != nil {
		return Prompt{}, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return Prompt{
		SystemPrompt:    prpt.SystemPrompt,
		StartUserPrompt: tplbuff.String(),
	}, nil
}
