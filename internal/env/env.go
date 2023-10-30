package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	TimeZone        string `json:"time_zone" env:"TZ" envDefault:"Etc/UTC" validate:"timezone"`
	Port            int    `json:"port" env:"PORT" envDefault:"8080" validate:"min=1024,max=49151"`
	LogLevel        string `json:"log_level" env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
	HealthcheckURL  string `json:"healthcheck_url" env:"HEALTHCHECK_URL" validate:"omitempty,url,endswith=/"`
	HealthcheckUUID string `json:"healthcheck_uuid" env:"HEALTHCHECK_UUID" validate:"omitempty,uuid"`
	BackupCron      string `json:"backup_cron" env:"BACKUP_CRON" validate:"omitempty,cron"`
	CleanupCron     string `json:"cleanup_cron" env:"CLEANUP_CRON" validate:"omitempty,cron"`
	CheckCron       string `json:"check_cron" env:"CHECK_CRON" validate:"omitempty,cron"`
	NtfyEndpoint    string `json:"ntfy_endpoint" env:"NTFY_ENDPOINT" validate:"omitempty,url,endswith=/"`
	NtfyToken       string `json:"ntfy_token" env:"NTFY_TOKEN,unset" validate:"omitempty"`
	NtfyTopic       string `json:"ntfy_topic" env:"NTFY_TOPIC" validate:"omitempty"`
	Version         string `json:"version" env:"APP_VERSION" envDefault:"v0.0.0"`
	Identifier      string `json:"identifier" env:"IDENTIFIER" envDefault:"GoBackup"`
	SwaggerHost     string `json:"swagger_host" env:"SWAGGER_HOST" validate:"omitempty,url"`
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
	setAllDefaultEnvs(cfg)
	return cfg, nil
}

func NewEnvValidator() *validator.Validate {
	return validator.New()
}

func validateContent(cfg *Config) error {
	validate := NewEnvValidator()
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

func setAllDefaultEnvs(cfg *Config) {
	os.Setenv("TZ", cfg.TimeZone)
	os.Setenv("PORT", fmt.Sprintf("%d", cfg.Port))
	os.Setenv("LOG_LEVEL", cfg.LogLevel)
	os.Setenv("IDENTIFIER", cfg.Identifier)
}
