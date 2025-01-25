# Usage

## Run Standard
You decide what GitHub issue you want to resolve

Example:

- GitHub working repository `clover0/example-repository`
- GitHub issue number 123 to solve
- Base branch `main` to create a pull request
- LLM is Anthropic Claude 3.5 Sonnet

```shell
$ issue-agent issue \
  --github_owner clover0 \
  --work_repository example-repository \
  --github_issue_number 123 \
  --base_branch main \
  --model claude-3-5-sonnet-latest \
  --language English \
  --log_level debug
```

## Run with Environment Variables

With environment variables in one line. [`gh` CLI is useful](https://github.com/cli/cli#installation).
```shell
$ GITHUB_TOKEN=$(gh auth token) \
  ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY} \
  issue-agent issue \
    --github_owner your_owner \
    --work_repository your_repository \
    --github_issue_number issue_number \
    --base_branch your_branch  \
    --model claude-3-5-sonnet-latest \
    --language English
```

OPENAI_API_KEY or ANTHROPIC_API_KEY environment variable is required


## Run AWS Bedrock with SSO session

```sh
$ GITHUB_TOKEN=$(gh auth token) \
issue-agent issue \
  --github_owner your_owner \
  --work_repository your_repository \
  --github_issue_number issue_number \
  --base_branch your_branch  \
  --model anthropic.claude-3-5-sonnet-20241022-v2:0 \
  --language Japanese \
  --log_level debug \
  --aws_profile your_profile \
  --aws_region us-east-1
```


## Branch

Working branch is created automatically. `agent-` is added to the branch prefix.
