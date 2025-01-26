# Setup

## Environment Variables
```shell

GITHUB_TOKEN=your github token

# If you use OpenAI models
OPENAI_API_KEY=your OpenAI API Key

# If you use Anthropic models
ANTHROPIC_API_KEY=your Anthropic API Key
```

##  More Configuration
Copy [default_config.yml](https://github.com/clover0/issue-agent/blob/main/agent/config/default_config.yml) to your repository root as `issue_agent.yml`
and edit the file.

If you want to create `issue_agent.yml` file in an arbitrary path,
specify the path in the command flags when executing the agent.

```shell 

issue-agent \ 
    --config_path /path/to/issue_agent.yml
```
