package env

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type Config struct {
	TimeZone             string `env:"TZ" envDefault:"Etc/UTC" validate:"timezone"`
	Port                 int    `env:"PORT" envDefault:"8156" validate:"omitempty,numeric"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
	NtfyUrl              string `env:"NTFY_URL" envDefault:"https://ntfy.sh/" validate:"omitempty,url,endswith=/"`
	NtfyTopic            string `env:"NTFY_TOPIC" envDefault:"gocron"`
	NtfyToken            string `env:"NTFY_TOKEN,unset"`
	SendMessageOnSuccess bool   `env:"SEND_ON_SUCCESS" envDefault:"true" validate:"omitempty,boolean"`
}

var errParse = errors.New("error parsing environment variables")

var logLevels = map[string]log.Lvl{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}

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

func (cfg *Config) GetLogLevel() log.Lvl {
	level := logLevels[cfg.LogLevel]
	return level
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
