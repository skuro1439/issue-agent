## Configuration YAML
At your repository root, create a `issue_agent.yml` file with the following content.

```yaml
# Communication language
# English, Japanese...
# Default is English
communication_language: "Japanese"

# Required
workdir: "/tmp/repositories/github-issue-agent"

# Required
# debug, info, error
log_level: "debug"

agent:
  # Default prompt template is embed config.
  # prompt_path is a relative path from the execution directory.
  # e.g) config/prompt_template_en.yml
  prompt_path: ""

  # Default prompt template is embed config.
  #  prompt_path: "prompt_en.yml"

  # LLM model name
  model: "gpt-4o"
  #  model: "claude-3-5-sonnet-20241022"

  # Maximum steps to run agent
  # The following are defined as 1 step
  # - user to LLM and returned to user from LLM
  # - execution function
  max_steps: 70

  skip_review_agents: true

  git:
    # Required
    user_name: "t.koenuma2@gmail.com"

    # Required
    user_email: "takeshi.koenuma"

  # GitHub environment for agent
  github:
    # Don't submit files to GitHub by Pull Request.
    no_submit: false

    # Whether to clone repository to the workdir
    clone_repository: true

    # Required
    # Repositories owner to operate
    #    owner: "clover0"
    owner: "reiwa5"

  # Allow agent to use function.
  # Belows are the default functions.
  allow_functions:
    - get_pull_request_diff
    #- get_web_page_from_url
    #- get_web_search_result
    - list_files
    - modify_file
    - open_file
    - put_file
    - submit_files
    - search_files
```
