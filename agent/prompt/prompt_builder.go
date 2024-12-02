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

func BuildRequirePrompt(promptTpl PromptTemplate, issueLoader loader.Loader, issueNumber string) (Prompt, error) {
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to load issue: %w", err)
	}

	return BuildPrompt(promptTpl, "requirement", map[string]any{
		"issue":       iss.Content,
		"issueNumber": issueNumber,
	})
}

func BuildDeveloperPrompt(promptTpl PromptTemplate, issueLoader loader.Loader, issueNumber string) (Prompt, error) {
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to load issue: %w", err)
	}

	return BuildPrompt(promptTpl, "developer", map[string]any{
		"issue":       iss.Content,
		"issueNumber": issueNumber,
	})
}

func BuildSecurityPrompt(promptTpl PromptTemplate, changedFilesPath []string) (Prompt, error) {
	return BuildPrompt(promptTpl, "security", map[string]any{
		"filePaths": changedFilesPath,
	})
}

func BuildPrompt(promptTpl PromptTemplate, templateName string, templateMap map[string]any) (Prompt, error) {
	var prpt Prompt
	for _, p := range promptTpl.Agents {
		if p.Name == templateName {
			prpt = Prompt{
				SystemPrompt:    p.SystemTemplate,
				StartUserPrompt: p.UserTemplate,
			}
			break
		}
	}

	if prpt.StartUserPrompt == "" {
		return Prompt{}, fmt.Errorf("failed to find %s prompt. You must have  name=%s prompt in the prompt template", templateName, templateName)
	}

	tpl, err := template.New("prompt").Parse(prpt.StartUserPrompt)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(tplbuff, templateMap); err != nil {
		return Prompt{}, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return Prompt{
		SystemPrompt:    prpt.SystemPrompt,
		StartUserPrompt: tplbuff.String(),
	}, nil
}
