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

func BuildRequirementPrompt(promptTpl PromptTemplate, issue loader.Issue) (Prompt, error) {
	return BuildPrompt(promptTpl, "requirement", map[string]any{
		"issue":       issue.Content,
		"issueNumber": issue.Path,
	})
}

func BuildDeveloperPrompt(promptTpl PromptTemplate, issueLoader loader.Loader, issueNumber string, instruction string) (Prompt, error) {
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to load issue: %w", err)
	}

	return BuildPrompt(promptTpl, "developer", map[string]any{
		"issue":       iss.Content,
		"issueNumber": issueNumber,
		"instruction": instruction,
	})
}

func BuildSecurityPrompt(promptTpl PromptTemplate, changedFilesPath []string) (Prompt, error) {
	m := make(map[string]any)

	m["filePaths"] = changedFilesPath

	m["noFiles"] = ""
	if len(changedFilesPath) == 0 {
		m["noFiles"] = "no changed files"
	}

	return BuildPrompt(promptTpl, "security", m)
}

func BuildReviewManagerPrompt(promptTpl PromptTemplate, issue loader.Issue, changedFilesPath []string) (Prompt, error) {
	m := make(map[string]any)

	m["filePaths"] = changedFilesPath
	m["issue"] = issue.Content

	m["noFiles"] = ""
	if len(changedFilesPath) == 0 {
		m["noFiles"] = "no changed files"
	}

	return BuildPrompt(promptTpl, "review-manager", m)
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
		return Prompt{}, fmt.Errorf("failed to find %s prompt. you must have  name=%s prompt in the prompt template", templateName, templateName)
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
