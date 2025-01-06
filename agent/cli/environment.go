package cli

// Environment variable names
const (
	AnthropicApiKey = "ANTHROPIC_API_KEY"
	GithubToken     = "GITHUB_TOKEN"
	OpenaiApiKey    = "OPENAI_API_KEY"
)

func EnvNames() []string {
	return []string{
		AnthropicApiKey,
		GithubToken,
		OpenaiApiKey,
	}
}
