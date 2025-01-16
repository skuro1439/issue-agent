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
	fs.StringVar(&cfg.Language, "language", "", "Language spoken agent")
	fs.StringVar(&cfg.Model, "model", "", "LLM Model name")
	fs.StringVar(&cfg.GitHubOwner, "github_owner", "", "GitHub owner of the repository")
}
