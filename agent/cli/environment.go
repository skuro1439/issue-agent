package cli

// Environment variable names
const (
	LOG_LEVEL       = "LOG_LEVEL"
	AnthropicApiKey = "ANTHROPIC_API_KEY"
	GithubToken     = "GITHUB_TOKEN"
	OpenaiApiKey    = "OPENAI_API_KEY"
)

func EnvNames() []string {
	return []string{
		LOG_LEVEL,
		AnthropicApiKey,
		GithubToken,
		OpenaiApiKey,
	}
}
