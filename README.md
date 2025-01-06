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
```
cp agent/config/default_config.yml agent/issue_agent.yml
```

```yaml
# Edit agent/config.yaml
# Example
workdir: "/tmp/repositories/github-issue-agent"
agent:
  prompt_template: ""
  model: "claude-3-5-sonnet-20241022"
  max_steps: 100
  git:
    user_name: "t.koenuma2@gmail.com"
    user_email: "takeshi.koenuma"
  github:
    no_submit: false
    clone_repository: true
    owner: "clover0"
    repository: "github-issue-agent"
    base_branch: "main"
  allow_functions:
    - get_pull_request_diff
    - get_web_page_from_url
    - get_web_search_result
    - list_files
    - modify_file
    - open_file
    - put_file
    - submit_file_service
    - submit_files
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
  -e LOG_LEVEL=debug \
  agent \
  go run cmd/runner/main.go issue \
    -config issue_agent.yml \
    -github_issue_number 123 \
    -base_branch master
```
  - Working branch is created automatically. (`agent-` prefix)
  - OPENAI_API_KEY or ANTHROPIC_API_KEY environment variable is required

- Human reviews of work product by agent
