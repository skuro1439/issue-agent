package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yml
var defaultConfig []byte

const (
	ConfigFilePath = "/agent/config/config.yml"
	PromptFilePath = "/agent/config/prompt.yml"
	DefaultWorkDir = "/agent/repositories"

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
)

type Config struct {
	CommunicationLanguage string `yaml:"communication_language"`
	WorkDir               string `yaml:"workdir"`
	LogLevel              string `yaml:"log_level" validate:"log_level"`
	Agent                 struct {
		PromptPath       string `yaml:"prompt_path"`
		Model            string `yaml:"model" validate:"required"`
		MaxSteps         int    `yaml:"max_steps" validate:"gte=0"`
		SkipReviewAgents *bool  `yaml:"skip_review_agents"`
		Git              struct {
			UserName  string `yaml:"user_name" validate:"required"`
			UserEmail string `yaml:"user_email" validate:"required"`
		} `yaml:"git"`
		GitHub struct {
			NoSubmit        *bool  `yaml:"no_submit"`
			CloneRepository *bool  `yaml:"clone_repository"`
			Owner           string `yaml:"owner" validate:"required"`
		}
		AllowFunctions []string `yaml:"allow_functions" validate:"required"`
	} `yaml:"agent" validate:"required"`
}

func isValidLogLevel(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, level := range []string{LogDebug, LogInfo, LogError} {
		if level == value {
			return true
		}
	}
	return false
}

func LoadDefault() (Config, error) {
	return Load(ConfigFilePath)
}

func Load(path string) (Config, error) {
	var cnfg Config

	var data []byte
	if path == "" {
		data = defaultConfig
	} else {
		file, err := os.Open(path)
		if err != nil {
			return cnfg, err
		}
		defer file.Close()

		data, err = io.ReadAll(file)
		if err != nil {
			return cnfg, err
		}
	}

	if err := yaml.Unmarshal(data, &cnfg); err != nil {
		return cnfg, err
	}

	cnfg = setDefaults(cnfg)

	if err := ValidateConfig(cnfg); err != nil {
		return cnfg, err
	}

	return cnfg, nil
}

func ValidateConfig(config Config) error {
	validate := validator.New()
	if err := validate.RegisterValidation("log_level", isValidLogLevel); err != nil {
		return err
	}
	if err := validate.Struct(config); err != nil {
		errs := err.(validator.ValidationErrors)
		return fmt.Errorf("validation failed: %w\n", errs)
	}
	return nil
}

func setDefaults(conf Config) Config {
	if conf.CommunicationLanguage == "" {
		conf.CommunicationLanguage = "English"
	}

	if conf.WorkDir == "" {
		conf.WorkDir = DefaultWorkDir
	}

	if conf.Agent.GitHub.NoSubmit == nil {
		noSubmit := false
		conf.Agent.GitHub.NoSubmit = &noSubmit
	}

	if conf.Agent.GitHub.CloneRepository == nil {
		clone := true
		conf.Agent.GitHub.CloneRepository = &clone
	}

	if conf.Agent.SkipReviewAgents == nil {
		skip := false
		conf.Agent.SkipReviewAgents = &skip
	}

	return conf
}
