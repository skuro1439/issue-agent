# Command

```shell
$ issue-agent help

Usage
  issue-agent <command> [flags]
Commands  help: Show usage of commands and flags
  help: Show usage of commands and flags
  version: Show version of issue-agent CLI
  issue:
    --aws_profile
        AWS profile to use for credentials
    --aws_region
        AWS region to use for credentials and Bedrock. Default is aws profile's session region
    --base_branch
        Base Branch for pull request
    --config
        Path to the configuration file. Default is `agent/config/default_config.yml in this project`
    --from_file
        Issue content from file path
    --github_issue_number
        GitHub issue number to solve
    --github_owner
        The GitHub account owner of the repository. Required
    --language
        Language spoken by agent. Default is English
    --log_level
        Log level. Default is `info`. If you want to see LLM completions, set it to `debug`
    --model
        LLM Model name. Default is `claude-3-5-sonnet-latest`
    --work_repository
        Working repository to develop and create pull request
```
