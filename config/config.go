package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	ConfigFolder = "./config/"
)

type Env struct {
	Key   string
	Value string
}

type Job struct {
	Name     string
	Cron     string
	Envs     []Env
	Commands []string
}

type JobDefaults struct {
	Cron         string
	Envs         []Env
	PreCommands  []string
	PostCommands []string
}

type HealthCheck struct {
	Authorization string `yaml:"authorization"`
	Type          string `validate:"omitempty,oneof=HEAD GET POST" yaml:"type"`
	Start         Url    `yaml:"start"`
	End           Url    `yaml:"end"`
	Failure       Url    `yaml:"failure"`
}

type Url struct {
	Url    string            `yaml:"url"`
	Params map[string]string `yaml:"params"`
	Body   string            `yaml:"body"`
}

func init() {
	os.Mkdir(ConfigFolder, os.ModePerm)
}

func New() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("time_zone", "Etc/UTC")
	viper.SetDefault("delete_runs_after_days", 7)

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8156)

	viper.SetDefault("healthcheck.type", "POST")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigFolder)

	viper.SetEnvPrefix("GC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		err = viper.WriteConfigAs(ConfigFolder + "config.yaml")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	os.Setenv("TZ", viper.GetString("time_zone"))
}

func ConfigLoaded() bool {
	return viper.ConfigFileUsed() != ""
}

func GetLogLevel() slog.Level {
	switch strings.ToLower(viper.GetString("log_level")) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func GetJobs() []Job {
	jobs := []Job{}
	viper.UnmarshalKey("jobs", &jobs)
	return jobs
}

func GetJobByName(name string) *Job {
	jobs := GetJobs()
	for _, job := range jobs {
		if strings.EqualFold(job.Name, strings.ToLower(name)) {
			return &job
		}
	}
	return nil
}

type OrderedEnvs struct {
	Order []string
	Data  map[string]string
}

func GetEnvsByJobName(name string) OrderedEnvs {
	data := make(map[string]string)
	order := []string{}

	addEnv := func(key, value string) {
		if _, exists := data[key]; !exists {
			order = append(order, key) // preserve first occurrence
		}
		data[key] = value // always overwrite value
	}

	job := GetJobByName(name)
	if job == nil {
		return OrderedEnvs{Order: order, Data: data}
	}

	// Defaults
	defaultEnvs := []Env{}
	viper.UnmarshalKey("job_defaults.envs", &defaultEnvs)
	for _, env := range defaultEnvs {
		addEnv(env.Key, env.Value)
	}

	// Job-specific
	for _, env := range job.Envs {
		addEnv(env.Key, env.Value)
	}

	return OrderedEnvs{Order: order, Data: data}
}

func GetCommandsByJobName(name string) []string {
	commands := []string{}

	job := GetJobByName(name)
	if job == nil {
		return commands
	}

	commands = append(commands, viper.GetStringSlice("job_defaults.pre_commands")...)
	commands = append(commands, job.Commands...)
	commands = append(commands, viper.GetStringSlice("job_defaults.post_commands")...)

	return commands
}

func GetHealthcheck() HealthCheck {
	var healthcheck HealthCheck
	viper.UnmarshalKey("healthcheck", &healthcheck)
	return healthcheck
}

func GetDeleteRunsAfterDays() int {
	return viper.GetInt("delete_runs_after_days")
}

func GetServer() string {
	return fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
}

func GetJobsCron(job *Job) string {
	cron := job.Cron
	if cron == "" {
		cron = viper.GetString("job_defaults.cron")
	}
	return cron
}
