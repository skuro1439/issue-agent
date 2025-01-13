# GitHub Issue Agent
Agent that, given an assignment, produces work products in the git repository.


## Models
The following models are supported.

- OpenAI models
  - gpt-4o
  - gpt-4o-mini
- Anthropic models
  - claude-3-5-sonnet


## Work Product
Only GitHub pull requests are supported.


## Usage
### Startup Example
- Set up the config file

Copy [default_config.yml](agent/config/default_config.yml) to your repository root as `issue_agent.yml`.

```shell


```yaml
# Example
# issue_agent.yml

communication_language: "Japanese"
agent:
  model: "claude-3-5-sonnet-20241022"
  max_steps: 70
  git:
    user_name: "username"
    user_email: "email@example.com"
  github:
    owner: "clover0"
    repository: "github-issue-agent"
```

- Human decides what issue they want to resolve
  - e.g)
  - GitHub Repository `clover0/example-repository`
  - GitHub Issue Number 123

- Run Agent with parameters below run example
```shell
cd agent

docker compose run --rm \
  -e GITHUB_TOKEN=$(gh auth token) \
  -e OPENAI_API_KEY=${OPENAI_API_KEY} \
  -e ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY} \
  agent \
  go run cmd/runner/main.go issue \
    -config issue_agent.yml \
    -github_issue_number 123 \
    -base_branch master
```
  - Working branch is created automatically. (`agent-` prefix)
  - OPENAI_API_KEY or ANTHROPIC_API_KEY environment variable is required

- Human reviews of work product by agent
