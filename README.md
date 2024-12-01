# Software Task Agent
Agents that, given an assignment, produce work products in the repository.


## Models
Only OpenAI models are supported.


## Work Product
Only pull requests to GitHub are supported.


## Usage
```shell
GITHUB_TOKEN="your GitHub Token" \
OPENAI_API_KEY="your OpenAI API key" \
LOG_LEVEL=debug \
go run cmd/runner/main.go \
  -template {propmpt template path} \
  -github_issue_number {GitHub issue number} \
  -repository_owner {repository owner} \
  -repository {repository name} \
  --model {model version} \
  -base_branch {base branch} \
  -workdir {your workdir}
```


### Example run

1. Human decides what issue they want to resolve
  - e.g)
  - At `example-repository` repository
  - GitHub Issue Number 123 in that GitHub repository
  - Local repository is `${HOME}/examples/example-repository`
1. Run Agent with parameters below run example
```shell
GITHUB_TOKEN="your GitHub Token" \
OPENAI_API_KEY="your OpenAI API key" \
LOG_LEVEL=debug \
go run cmd/runner/main.go \
  -template ./agent/config/template/default_prompt_ja.yaml \
  -github_issue_number 123 \
  -repository_owner clover0 \
  -repository example-repository \
  --model gpt-4o \
  -base_branch master \
  -workdir ${HOME}/examples/example-repository
```
  - Working branch is created automatically
1. Human review of work product by agent
