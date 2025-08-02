package config

import (
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const (
	configFolder = "./config/"
)

var logLevels = map[string]log.Lvl{
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"off":   5,
}

func init() {
	os.Mkdir(configFolder, os.ModePerm)
}

func New() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("time_zone", "Etc/UTC")
	viper.SetDefault("delete_runs_after_days", 7)

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFolder)

	viper.SetEnvPrefix("GC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()

	os.Setenv("TZ", viper.GetString("time_zone"))
}

func ConfigLoaded() bool {
	return viper.ConfigFileUsed() != ""
}

func GetLogLevel() log.Lvl {
	if !ConfigLoaded() {
		return log.INFO
	}
	level := logLevels[viper.GetString("log_level")]
	return level
}
