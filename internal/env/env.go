package env

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	TimeZone  string `env:"TZ" envDefault:"Etc/UTC" validate:"timezone"`
	LogLevel  string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
	NtfyUrl   string `env:"NTFY_URL" envDefault:"https://ntfy.sh/" validate:"omitempty,url,endswith=/"`
	NtfyTopic string `env:"NTFY_TOPIC" envDefault:"gocron"`
	NtfyToken string `env:"NTFY_TOKEN,unset"`
}

var errParse = errors.New("error parsing environment variables")

func Parse() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}
	if err := validateContent(cfg); err != nil {
		return cfg, err
	}
	setTZDefaultEnv(cfg)
	return cfg, nil
}

func validateContent(cfg *Config) error {
	validate := validator.New()
	err := validate.Struct(cfg)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		} else {
			for _, err := range err.(validator.ValidationErrors) {
				return err
			}
		}
		return errParse
	}
	return nil
}

func setTZDefaultEnv(e *Config) {
	os.Setenv("TZ", e.TimeZone)
}
