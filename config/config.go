package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/flohoss/gocron/pkg/expand"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "./config/config.yaml"
)

var cfg GlobalConfig
var configFile = defaultConfigFile

var validate *validator.Validate
var mu sync.RWMutex

type GlobalConfig struct {
	LogLevel            string           `mapstructure:"log_level" validate:"omitempty,oneof=debug info warn error"`
	TimeZone            string           `mapstructure:"time_zone" validate:"omitempty,timezone"`
	DeleteRunsAfterDays int              `mapstructure:"delete_runs_after_days" validate:"gte=0"`
	Jobs                []Job            `mapstructure:"jobs" validate:"omitempty,dive"`
	JobDefaults         JobDefaults      `mapstructure:"job_defaults"`
	Healthcheck         HealthCheck      `mapstructure:"healthcheck" validate:"omitempty"`
	Server              ServerSettings   `mapstructure:"server"`
	Terminal            TerminalSettings `mapstructure:"terminal" validate:"omitempty"`
	Software            []Software       `mapstructure:"software" validate:"omitempty,dive"`
}

type Software struct {
	Name    string `mapstructure:"name" validate:"required"`
	Version string `mapstructure:"version"`
}

type ServerSettings struct {
	Address string `mapstructure:"address" validate:"required,ipv4"`
	Port    int    `mapstructure:"port" validate:"required,gte=1024,lte=65535"`
}

type Env struct {
	Key   string `mapstructure:"key" validate:"required"`
	Value string `mapstructure:"value" validate:"required"`
}

type Job struct {
	Name            string   `mapstructure:"name" validate:"required" json:"name"`
	Cron            string   `mapstructure:"cron" validate:"omitempty,cron" json:"cron"`
	DisableCron     bool     `mapstructure:"disable_cron" json:"disable_cron"`
	DisableFailFast bool     `mapstructure:"disable_fail_fast" json:"disable_fail_fast"`
	Envs            []Env    `mapstructure:"envs" validate:"dive" json:"-"`
	Commands        []string `mapstructure:"commands" validate:"required" json:"-"`
	Disabled        bool     `json:"disabled"`
}

type JobDefaults struct {
	Cron         string   `mapstructure:"cron" validate:"omitempty,cron"`
	Envs         []Env    `mapstructure:"envs" validate:"dive"`
	PreCommands  []string `mapstructure:"pre_commands"`
	PostCommands []string `mapstructure:"post_commands"`
}

type HealthCheck struct {
	Authorization string `mapstructure:"authorization"`
	Type          string `mapstructure:"type" validate:"omitempty,oneof=HEAD GET POST"`
	Start         Url    `mapstructure:"start" validate:"omitempty"`
	End           Url    `mapstructure:"end" validate:"omitempty"`
	Failure       Url    `mapstructure:"failure" validate:"omitempty"`
}

type Url struct {
	Url    string         `mapstructure:"url" validate:"required,url"`
	Params map[string]any `mapstructure:"params"`
	Body   string         `mapstructure:"body"`
}

type AllowedCommands struct {
	AllowAllArgs bool     `mapstructure:"allow_all_args"`
	Args         []string `mapstructure:"args"`
	// New field to store pre-processed arguments
	AllowedArgsMap map[string]struct{}
}

type TerminalSettings struct {
	AllowAllCommands bool                       `mapstructure:"allow_all_commands"`
	AllowedCommands  map[string]AllowedCommands `mapstructure:"allowed_commands" validate:"required_if=AllowAllCommands false,dive"`
}

func init() {
	validate = validator.New()
}

func defaultStarterJobs() []Job {
	return []Job{
		{
			Name:            "Example Scheduled Happy Path",
			Cron:            "0 5 * * 0",
			DisableFailFast: false,
			Commands: []string{
				"echo \"start\"",
				"date",
				"echo \"done\"",
			},
		},
		{
			Name:            "Example Continue On Failure",
			Cron:            "15 5 * * 0",
			DisableFailFast: true,
			Commands: []string{
				"echo \"before fail\"",
				"false",
				"echo \"continues\"",
			},
		},
		{
			Name: "Example Env Expansion",
			Cron: "30 5 * * 0",
			Envs: []Env{
				{Key: "TEST_VALUE", Value: "default-value"},
			},
			Commands: []string{
				"echo \"value=${TEST_VALUE}\"",
			},
		},
		{
			Name:        "Example Manual Long Running",
			DisableCron: true,
			Commands: []string{
				"sleep 2",
				"echo \"manual done\"",
			},
		},
	}
}

func New(configFilePath string) {
	SetConfigFilePath(configFilePath)
	configFolder := GetConfigFolderPath()

	if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
		slog.Error("Failed to create configuration directory", "error", err)
		os.Exit(1)
	}

	viper.SetDefault("log_level", "info")
	viper.SetDefault("time_zone", "UTC")
	viper.SetDefault("delete_runs_after_days", 7)
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8156)
	viper.SetDefault("healthcheck.type", "POST")
	viper.SetDefault("terminal.allow_all_commands", false)
	viper.SetDefault("jobs", defaultStarterJobs())

	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("GC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if notFound || os.IsNotExist(err) {
			err = viper.WriteConfigAs(configFile)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		} else {
			slog.Error("Failed to read configuration file", "error", err)
			os.Exit(1)
		}
	}

	viper.AutomaticEnv()

	if err := ValidateAndLoadConfig(viper.GetViper()); err != nil {
		slog.Error("Initial configuration validation failed", "error", err)
		os.Exit(1)
	}
}

func ValidateAndLoadConfig(v *viper.Viper) error {
	var tempCfg GlobalConfig
	if err := v.Unmarshal(&tempCfg); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	expand.ExpandEnvStrings(&tempCfg.Healthcheck)
	tempCfg.Terminal.Hydrate()

	if err := validate.Struct(tempCfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	mu.Lock()
	cfg = tempCfg
	mu.Unlock()

	os.Setenv("TZ", cfg.TimeZone)
	return nil
}

func ConfigLoaded() bool {
	return viper.ConfigFileUsed() != ""
}

func GetDefaultConfigFolder() string {
	return filepath.Dir(defaultConfigFile)
}

func GetDefaultConfigFile() string {
	return defaultConfigFile
}

func SetConfigFolderPath(folder string) {
	if folder == "" {
		folder = GetDefaultConfigFolder()
	}

	SetConfigFilePath(filepath.Join(folder, "config.yaml"))
}

func SetConfigFilePath(file string) {
	if file == "" {
		file = defaultConfigFile
	}

	resolvedFile := filepath.Clean(file)
	mu.Lock()
	configFile = resolvedFile
	mu.Unlock()
}

func GetConfigFilePath() string {
	mu.RLock()
	defer mu.RUnlock()
	return configFile
}

func GetConfigFolderPath() string {
	mu.RLock()
	defer mu.RUnlock()
	return configFolder
}

func GetLogLevel() slog.Level {
	mu.RLock()
	defer mu.RUnlock()
	switch strings.ToLower(cfg.LogLevel) {
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
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Jobs
}

func GetJobByName(name string) *Job {
	mu.RLock()
	defer mu.RUnlock()
	for _, job := range cfg.Jobs {
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
	mu.RLock()
	defer mu.RUnlock()
	data := make(map[string]string)
	order := []string{}

	addEnv := func(key, value string) {
		if _, exists := data[key]; !exists {
			order = append(order, key)
		}
		data[key] = value
	}

	job := GetJobByName(name)
	if job == nil {
		return OrderedEnvs{Order: order, Data: data}
	}

	for _, env := range cfg.JobDefaults.Envs {
		addEnv(env.Key, env.Value)
	}

	for _, env := range job.Envs {
		addEnv(env.Key, env.Value)
	}

	return OrderedEnvs{Order: order, Data: data}
}

func GetCommandsByJobName(name string) []string {
	mu.RLock()
	defer mu.RUnlock()
	commands := []string{}

	job := GetJobByName(name)
	if job == nil {
		return commands
	}

	commands = append(commands, cfg.JobDefaults.PreCommands...)
	commands = append(commands, job.Commands...)
	commands = append(commands, cfg.JobDefaults.PostCommands...)

	return commands
}

func GetHealthcheck() HealthCheck {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Healthcheck
}

func GetDeleteRunsAfterDays() int {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.DeleteRunsAfterDays
}

func GetServer() string {
	mu.RLock()
	defer mu.RUnlock()
	return fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
}

func GetJobsCron(job *Job) string {
	mu.RLock()
	defer mu.RUnlock()
	cron := job.Cron
	if cron == "" {
		cron = cfg.JobDefaults.Cron
	}
	return cron
}

func GetAllCrons() map[string][]Job {
	mu.RLock()
	defer mu.RUnlock()
	var cronJobs = make(map[string][]Job)
	jobs := GetJobs()

	for _, job := range jobs {
		if job.DisableCron {
			continue
		}
		cron := GetJobsCron(&job)
		cronJobs[cron] = append(cronJobs[cron], job)
	}

	return cronJobs
}

func GetTerminalSettings() TerminalSettings {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Terminal
}

func (s *TerminalSettings) Hydrate() {
	for cmdName, cmdConfig := range s.AllowedCommands {
		if len(cmdConfig.Args) > 0 {
			cmdConfig.AllowedArgsMap = make(map[string]struct{}, len(cmdConfig.Args))
			for _, arg := range cmdConfig.Args {
				cmdConfig.AllowedArgsMap[arg] = struct{}{}
			}
			s.AllowedCommands[cmdName] = cmdConfig
		}
	}
}

func EnableAllJobs() {
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.Jobs {
		cfg.Jobs[i].Disabled = false
	}
}

func DisableAllJobs() {
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.Jobs {
		cfg.Jobs[i].Disabled = true
	}
}

func EnableScheduledJobs() {
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.Jobs {
		cfg.Jobs[i].Disabled = cfg.Jobs[i].DisableCron
	}
}

func EnableNonScheduledJobs() {
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.Jobs {
		cfg.Jobs[i].Disabled = !cfg.Jobs[i].DisableCron
	}
}

func ToggleDisabledJob(name string) error {
	mu.Lock()
	defer mu.Unlock()
	for i, job := range cfg.Jobs {
		if strings.EqualFold(job.Name, name) {
			cfg.Jobs[i].Disabled = !cfg.Jobs[i].Disabled
			return nil
		}
	}
	return fmt.Errorf("job %q not found", name)
}
