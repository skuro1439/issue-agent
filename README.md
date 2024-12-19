# GitHub Issue Agent
Agent that, given an assignment, produces work products in the repository.


## Models
- OpenAI models
  - gpt-4o
  - gpt-4o-mini
- Anthropic models
  - claude-3-5-sonnet

are supported.


## Work Product
Only GitHub pull requests are supported.


## Usage
### Example run

1. Human decides what issue they want to resolve
  - e.g)
  - At `example-repository` repository
  - GitHub Issue Number 123 in that GitHub repository
  - GitHub Repository `clover0/example-repository`
1. Run Agent with parameters below run example
```shell
cd agent
docker compose run --rm \
  -e GITHUB_TOKEN=$(gh auth token) \
  -e OPENAI_API_KEY=${OPENAI_API_KEY} \
  -e ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY} \
  -e LOG_LEVEL=debug \
  agent \
  go run cmd/runner/main.go issue \
    -github_issue_number 123 \
    -clone_repository \
    -repository_owner clover0 \
    -repository example-repository \
    -model gpt-4o \
    -base_branch master \
    -workdir /usr/local/repositories/example-repository \
    -git_email email@example.com \
    -git_name clover0
```
  - Working branch is created automatically
  - Git clone at /usr/local/repositories
  - LLM OPENAI_API_KEY or ANTHROPIC_API_KEY is required
1. Human review of work product by agent
