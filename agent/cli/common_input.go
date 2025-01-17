package cli

import (
	"flag"
)

type CommonInput struct {
	Config      string
	Language    string
	Model       string
	GitHubOwner string
}

func addCommonFlags(fs *flag.FlagSet, cfg *CommonInput) {
	fs.StringVar(&cfg.Config, "config", "", "Path to the configuration file. Default is `agent/config/default_config.yml in this project`")
	fs.StringVar(&cfg.Language, "language", "", "Language spoken by agent. Default is English")
	fs.StringVar(&cfg.Model, "model", "", "LLM Model name. Default is `claude-3-5-sonnet-latest`")
	fs.StringVar(&cfg.GitHubOwner, "github_owner", "", "The GitHub account owner of the repository. Required")
}
