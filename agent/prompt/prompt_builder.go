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

func BuildDeveloperPrompt(promptTpl PromptTemplate, issueLoader loader.Loader, issueNumber string) (Prompt, error) {
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		panic(err)
	}

	m := map[string]string{
		"issue":       iss.Content,
		"issueNumber": issueNumber,
	}

	var prpt Prompt
	for _, p := range promptTpl.Agents {
		if p.Name == "developer" {
			prpt = Prompt{
				SystemPrompt:    p.SystemTemplate,
				StartUserPrompt: p.UserTemplate,
			}
			break
		}
	}

	if prpt.StartUserPrompt == "" {
		return Prompt{}, fmt.Errorf("failed to find developer prompt. You must have  name=developer prompt in the prompt template")
	}

	usrTpl, err := template.New("prompt").Parse(prpt.StartUserPrompt)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := usrTpl.Execute(tplbuff, m); err != nil {
		return Prompt{}, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return Prompt{
		SystemPrompt:    prpt.SystemPrompt,
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

	if prpt.StartUserPrompt == "" {
		return Prompt{}, fmt.Errorf("failed to find security prompt. You must have  name=security prompt in the prompt template")
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
