package env

import (
	"fmt"
	"gobackup/validate"
	"log"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	TimeZone        string `env:"TZ" envDefault:"Etc/UTC"`
	Port            int    `env:"PORT" envDefault:"8080"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error dpanic panic fatal"`
	HealthcheckURL  string `env:"HEALTHCHECK_URL" validate:"omitempty,url,endswith=/"`
	HealthcheckUUID string `env:"HEALTHCHECK_UUID" validate:"omitempty,uuid"`
	BackupCron      string `env:"BACKUP_CRON" validate:"omitempty,cron"`
	CleanupCron     string `env:"CLEANUP_CRON" validate:"omitempty,cron"`
	CheckCron       string `env:"CHECK_CRON" validate:"omitempty,cron"`
	NotificationURL string `env:"NOTIFICATION_URL" validate:"omitempty,shoutrrr"`
	Version         string `env:"VERSION" envDefault:"v0.0.0"`
	Identifier      string `env:"IDENTIFIER" envDefault:"GoBackup"`
	DefaultSubset   uint   `env:"DEFAULT_SUBSET" envDefault:"10"`
}

func Parse(val *validator.Validate) *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalln(err.Error())
	}
	err := val.Struct(cfg)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println(validate.MessageForError(err))
		}
		os.Exit(1)
	}
	setAllDefaultEnvs(cfg)
	return cfg
}

func setAllDefaultEnvs(cfg *Config) {
	os.Setenv("TZ", cfg.TimeZone)
	os.Setenv("PORT", fmt.Sprintf("%d", cfg.Port))
	os.Setenv("LOG_LEVEL", cfg.LogLevel)
	os.Setenv("IDENTIFIER", cfg.Identifier)
	os.Setenv("DEFAULT_SUBSET", fmt.Sprintf("%d", cfg.DefaultSubset))
}
