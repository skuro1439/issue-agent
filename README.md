# GitHub Issue Agent
Agent that, given an assignment, produces work products in the repository.


## Models
Only OpenAI models are supported.


## Work Product
Only pull requests to GitHub are supported.


## Usage

### Example run

1. Human decides what issue they want to resolve
  - e.g)
  - At `example-repository` repository
  - GitHub Issue Number 123 in that GitHub repository
  - GitHub Repository `clover0/example-repository`
1. Run Agent with parameters below run example
```shell
docker compose run --rm \
  -e GITHUB_TOKEN=$(gh auth token) \
  -e OPENAI_API_KEY=${OPENAI_API_KEY} \
  -e LOG_LEVEL=debug \
  agent \
  go run cmd/runner/main.go \
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
1. Human review of work product by agent
