package env

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/containrrr/shoutrrr"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	TimeZone        string `env:"TZ" envDefault:"Etc/UTC" validate:"timezone"`
	Port            int    `env:"PORT" envDefault:"8080" validate:"number,min=1024,max=49151"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error panic fatal"`
	HealthcheckURL  string `env:"HEALTHCHECK_URL" validate:"omitempty,url,endswith=/"`
	HealthcheckUUID string `env:"HEALTHCHECK_UUID" validate:"omitempty,uuid"`
	BackupCron      string `env:"BACKUP_CRON" validate:"omitempty,cron"`
	CleanupCron     string `env:"CLEANUP_CRON" validate:"omitempty,cron"`
	CheckCron       string `env:"CHECK_CRON" validate:"omitempty,cron"`
	NotificationURL string `env:"NOTIFICATION_URL" validate:"omitempty,shoutrrr"`
	Version         string `env:"VERSION" envDefault:"v0.0.0"`
	Identifier      string `env:"IDENTIFIER" envDefault:"GoBackup"`
	DefaultSubset   uint   `env:"DEFAULT_SUBSET" envDefault:"10" validate:"number,min=1,max=100"`
}

func Parse() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalln(err)
	}
	if err := validateContent(cfg); err != nil {
		log.Fatalln(err)
	}
	setAllDefaultEnvs(cfg)
	return cfg
}

func newEnvValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation(`shoutrrr`, func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(string)
		_, err := shoutrrr.CreateSender(value)
		return err == nil
	})

	return validate
}

func validateContent(cfg *Config) error {
	validate := newEnvValidator()
	err := validate.Struct(cfg)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
		} else {
			for _, err := range err.(validator.ValidationErrors) {
				log.Println(err)
			}
		}
		return fmt.Errorf("error parsing environment variables")
	}
	return nil
}

func setAllDefaultEnvs(cfg *Config) {
	os.Setenv("TZ", cfg.TimeZone)
	os.Setenv("PORT", fmt.Sprintf("%d", cfg.Port))
	os.Setenv("LOG_LEVEL", cfg.LogLevel)
	os.Setenv("IDENTIFIER", cfg.Identifier)
	os.Setenv("DEFAULT_SUBSET", fmt.Sprintf("%d", cfg.DefaultSubset))
}
