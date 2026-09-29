package config

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/flohoss/gocron/internal/validate"
	"github.com/flohoss/gocron/pkg/expand"
	mapstructure "github.com/go-viper/mapstructure/v2"
	goslug "github.com/gosimple/slug"
	"github.com/spf13/viper"
)

const defaultConfigFile = "./config/config.yaml"

//go:embed config.example.yaml
var defaultConfig []byte

var cfg GlobalConfig
var configFile = defaultConfigFile
var logLevel slog.LevelVar

var mu sync.RWMutex

type GlobalConfig struct {
	LogLevel            slog.Level       `mapstructure:"log_level"`
	TimeZone            *time.Location   `mapstructure:"time_zone"`
	DeleteRunsAfterDays int              `mapstructure:"delete_runs_after_days" validate:"gte=0"`
	DB                  DBSettings       `mapstructure:"db"`
	Jobs                []Job            `mapstructure:"jobs" validate:"omitempty,dive"`
	JobDefaults         JobDefaults      `mapstructure:"job_defaults"`
	Healthcheck         HealthCheck      `mapstructure:"healthcheck" validate:"omitempty"`
	Server              ServerSettings   `mapstructure:"server"`
	Terminal            TerminalSettings `mapstructure:"terminal" validate:"omitempty"`
	Software            []Software       `mapstructure:"software" validate:"omitempty,dive"`
	Auth                AuthSettings     `mapstructure:"auth" validate:"omitempty"`
}

type DBSettings struct {
	Location string `mapstructure:"location"`
	Name     string `mapstructure:"name"`
}

type Software struct {
	Name    string `mapstructure:"name" validate:"required"`
	Version string `mapstructure:"version"`
}

type ServerSettings struct {
	Address        string            `mapstructure:"address" validate:"required,ipv4"`
	Port           int               `mapstructure:"port" validate:"required,gte=1024,lte=65535"`
	TrustedProxies []string          `mapstructure:"trusted_proxies" validate:"omitempty,dive,cidr"`
	CORS           CORSSettings      `mapstructure:"cors"`
	RateLimit      RateLimitSettings `mapstructure:"rate_limit"`
}

type CORSSettings struct {
	AllowOrigins []string `mapstructure:"allow_origins" validate:"omitempty,dive,omitempty,cors_origin"`
}

type RateLimitSettings struct {
	Enabled bool    `mapstructure:"enabled"`
	Rate    float64 `mapstructure:"rate" validate:"gte=0"`
	Burst   int     `mapstructure:"burst" validate:"gte=0"`
}

type Env struct {
	Key   string `mapstructure:"key" validate:"required,env_key"`
	Value string `mapstructure:"value" validate:"required"`
}

type Job struct {
	Name            string        `mapstructure:"name" validate:"required" json:"name"`
	Slug            string        `mapstructure:"-" json:"slug"`
	Cron            string        `mapstructure:"cron" validate:"omitempty,cron" json:"cron"`
	DisableCron     bool          `mapstructure:"disable_cron" json:"disable_cron"`
	DisableFailFast bool          `mapstructure:"disable_fail_fast" json:"disable_fail_fast"`
	Timeout         time.Duration `mapstructure:"timeout" validate:"gte=0" json:"timeout,omitempty"`
	Retries         int           `mapstructure:"retries" validate:"gte=0" json:"retries,omitempty"`
	Envs            []Env         `mapstructure:"envs" validate:"dive" json:"-"`
	Commands        []string      `mapstructure:"commands" validate:"required" json:"-"`
	Disabled        bool          `json:"disabled"`
}

type JobDefaults struct {
	Cron         string        `mapstructure:"cron" validate:"omitempty,cron"`
	Timeout      time.Duration `mapstructure:"timeout" validate:"gte=0"`
	Retries      int           `mapstructure:"retries" validate:"gte=0"`
	Envs         []Env         `mapstructure:"envs" validate:"dive"`
	PreCommands  []string      `mapstructure:"pre_commands"`
	PostCommands []string      `mapstructure:"post_commands"`
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

type AuthSettings struct {
	OIDC OIDCSettings `mapstructure:"oidc" validate:"omitempty"`
}

type OIDCSettings struct {
	Enabled      bool          `mapstructure:"enabled"`
	IssuerURL    string        `mapstructure:"issuer_url" validate:"required_if=Enabled true,omitempty,url,endsnotwith=/"`
	ClientID     string        `mapstructure:"client_id" validate:"required_if=Enabled true"`
	ClientSecret string        `mapstructure:"client_secret" validate:"required_if=Enabled true"`
	SessionTTL   time.Duration `mapstructure:"session_ttl" validate:"gte=0"`
	CookieSecure bool          `mapstructure:"cookie_secure"`
}

func slugifyJobName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}

	return goslug.Make(trimmed)
}

func New(configFilePath string) {
	SetConfigFilePath(configFilePath)
	configFolder := GetConfigFolderPath()

	if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
		slog.Error("Failed to create configuration directory", "error", err)
		os.Exit(1)
	}

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8156)
	viper.SetDefault("server.rate_limit.rate", 20)

	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("GC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if !notFound && !os.IsNotExist(err) {
			slog.Error("Failed to read configuration file", "error", err)
			os.Exit(1)
		}

		if err := os.WriteFile(configFile, defaultConfig, 0o644); err != nil {
			slog.Error("Failed to write example configuration", "error", err)
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			slog.Error("Failed to read written example configuration", "error", err)
			os.Exit(1)
		}
	}

	viper.AutomaticEnv()

	if err := ValidateAndLoadConfig(viper.GetViper()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ValidateAndLoadConfig(v *viper.Viper) error {
	var tempCfg GlobalConfig
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToTimeLocationHookFunc(),
		mapstructure.TextUnmarshallerHookFunc(),
	)
	if err := v.Unmarshal(&tempCfg, viper.DecodeHook(decodeHook)); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	expand.ExpandEnvStrings(&tempCfg.Healthcheck)
	tempCfg.Terminal.Hydrate()

	seenSlugs := make(map[string]string, len(tempCfg.Jobs))
	for i := range tempCfg.Jobs {
		tempCfg.Jobs[i].Slug = slugifyJobName(tempCfg.Jobs[i].Name)
		if tempCfg.Jobs[i].Slug == "" {
			return fmt.Errorf("configuration validation failed: jobs[%d].name must produce a non-empty slug", i)
		}
		if previous, exists := seenSlugs[tempCfg.Jobs[i].Slug]; exists {
			return fmt.Errorf("configuration validation failed: jobs[%d].name slug %q collides with job %q", i, tempCfg.Jobs[i].Slug, previous)
		}
		seenSlugs[tempCfg.Jobs[i].Slug] = tempCfg.Jobs[i].Name
	}

	if err := validate.Struct(tempCfg); err != nil {
		return fmt.Errorf("configuration validation failed:\n%s", err)
	}

	if tempCfg.Server.RateLimit.Enabled && tempCfg.Server.RateLimit.Rate <= 0 {
		return fmt.Errorf("configuration validation failed: server.rate_limit.rate must be greater than 0 when the rate limiter is enabled")
	}

	mu.Lock()
	cfg = tempCfg
	mu.Unlock()

	logLevel.Set(tempCfg.LogLevel)

	if tempCfg.TimeZone != nil {
		os.Setenv("TZ", tempCfg.TimeZone.String())
	}

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
	return filepath.Dir(GetConfigFilePath())
}

func GetDBLocation() string {
	mu.RLock()
	configPath := configFile
	dbLocation := cfg.DB.Location
	defer mu.RUnlock()

	baseDir := filepath.Dir(configPath)
	if dbLocation == "" {
		return baseDir
	}
	if filepath.IsAbs(dbLocation) {
		return filepath.Clean(dbLocation)
	}

	return filepath.Clean(filepath.Join(baseDir, dbLocation))
}

func GetDBName() string {
	mu.RLock()
	name := cfg.DB.Name
	defer mu.RUnlock()

	if name == "" {
		return "db.sqlite"
	}

	cleanName := filepath.Base(filepath.Clean(name))
	if cleanName == "." || cleanName == string(filepath.Separator) {
		return "db.sqlite"
	}

	return cleanName
}

func GetLogLevel() slog.Level {
	return logLevel.Level()
}

func LogLevelVar() *slog.LevelVar {
	return &logLevel
}

func GetLocation() *time.Location {
	mu.RLock()
	defer mu.RUnlock()

	if cfg.TimeZone == nil {
		return time.Local
	}

	return cfg.TimeZone
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
		if job.Slug == name || strings.EqualFold(job.Name, name) {
			return &job
		}
	}
	return nil
}

type OrderedEnvs struct {
	Order []string
	Data  map[string]string
}

func GetEnvsForJob(job *Job) OrderedEnvs {
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

	for _, env := range cfg.JobDefaults.Envs {
		addEnv(env.Key, env.Value)
	}

	for _, env := range job.Envs {
		addEnv(env.Key, env.Value)
	}

	return OrderedEnvs{Order: order, Data: data}
}

func GetCommandsForJob(job *Job) []string {
	mu.RLock()
	defer mu.RUnlock()
	commands := []string{}

	commands = append(commands, cfg.JobDefaults.PreCommands...)
	commands = append(commands, job.Commands...)
	commands = append(commands, cfg.JobDefaults.PostCommands...)

	return commands
}

func GetHealthcheck() HealthCheck {
	mu.RLock()
	hc := cfg.Healthcheck
	mu.RUnlock()

	if hc.Type == "" {
		hc.Type = "POST"
	}

	return hc
}

func GetSoftware() []Software {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Software
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

func GetCORSSettings() CORSSettings {
	mu.RLock()
	defer mu.RUnlock()

	settings := cfg.Server.CORS
	if len(settings.AllowOrigins) == 0 {
		settings.AllowOrigins = []string{"*"}
	}
	return settings
}

func GetRateLimitSettings() RateLimitSettings {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Server.RateLimit
}

func GetTrustedProxies() []string {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Server.TrustedProxies
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

func GetTimeoutForJob(job *Job) time.Duration {
	mu.RLock()
	defer mu.RUnlock()
	if job.Timeout > 0 {
		return job.Timeout
	}
	return cfg.JobDefaults.Timeout
}

func GetRetriesForJob(job *Job) int {
	mu.RLock()
	defer mu.RUnlock()
	if job.Retries > 0 {
		return job.Retries
	}
	return cfg.JobDefaults.Retries
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
		if cron == "" {
			continue
		}
		cronJobs[cron] = append(cronJobs[cron], job)
	}

	return cronJobs
}

func GetTerminalSettings() TerminalSettings {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Terminal
}

func GetAuth() AuthSettings {
	mu.RLock()
	defer mu.RUnlock()
	return cfg.Auth
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
