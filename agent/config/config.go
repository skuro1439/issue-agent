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

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
)

type Config struct {
	WorkDir  string `yaml:"workdir" validate:"required"`
	LogLevel string `yaml:"log_level" validate:"log_level"`
	Agent    struct {
		PromptTemplate string `yaml:"prompt_template"`
		Model          string `yaml:"model" validate:"required"`
		MaxSteps       int    `yaml:"max_steps" validate:"gte=0"`
		Git            struct {
			UserName  string `yaml:"user_name" validate:"required"`
			UserEmail string `yaml:"user_email" validate:"required"`
		} `yaml:"git"`
		GitHub struct {
			NoSubmit        bool   `yaml:"no_submit"`
			CloneRepository bool   `yaml:"clone_repository"`
			Owner           string `yaml:"owner" validate:"required"`
			Repository      string `yaml:"repository" validate:"required"`
			BaseBranch      string `yaml:"base_branch" validate:"required"`
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

func Load(path string) (Config, error) {
	var cnfg Config

	var data []byte
	if path == "" {
		data = defaultConfig
	} else {
		path = ConfigFilePath
		file, err := os.Open(path)
		if err != nil {
			return cnfg, err
		}
		data, err = io.ReadAll(file)
		if err != nil {
			return cnfg, err
		}
	}

	if err := yaml.Unmarshal(data, &cnfg); err != nil {
		return cnfg, err
	}

	return cnfg, ValidateConfig(cnfg)
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
