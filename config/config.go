package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Env struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Command struct {
	Command string `yaml:"command"`
}

type Job struct {
	Name     string    `yaml:"name"`
	Cron     string    `yaml:"cron,omitempty"`
	Envs     []Env     `yaml:"envs"`
	Commands []Command `yaml:"commands"`
}

type Config struct {
	Defaults struct {
		Cron string `yaml:"cron"`
		Envs []Env  `yaml:"envs"`
	} `yaml:"defaults"`
	Jobs []Job `yaml:"jobs"`
}

const (
	cronRegexString   = `(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|((\*|\d+)(\/|-)\d+)|\d+|\*) ?){5,7})`
	envKeyRegexString = `^[A-Z_][A-Z0-9_]*$`
)

func (c *Config) Validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("please specify at least one job")
	}
	for _, job := range c.Jobs {
		if job.Name == "" {
			return fmt.Errorf("please specify a name for each job")
		}
		if job.Cron == "" {
			return fmt.Errorf("please specify a cron for each job")
		} else {
			re := regexp.MustCompile(cronRegexString)
			if !re.MatchString(job.Cron) {
				return fmt.Errorf("please specify a valid cron for each job")
			}
		}
		if len(job.Commands) == 0 {
			return fmt.Errorf("please specify at least one command for each job")
		}
		for _, env := range job.Envs {
			if env.Key != "" {
				re := regexp.MustCompile(envKeyRegexString)
				if !re.MatchString(env.Key) {
					return fmt.Errorf("please specify a valid key for each environment variable")
				}
			}
		}
	}
	return nil
}

func mergeEnvs(defaultEnvs, jobEnvs []Env) []Env {
	envMap := make(map[string]string)
	for _, env := range defaultEnvs {
		envMap[env.Key] = env.Value
	}
	for _, env := range jobEnvs {
		envMap[env.Key] = env.Value
	}

	mergedEnvs := make([]Env, 0, len(envMap))
	for key, value := range envMap {
		mergedEnvs = append(mergedEnvs, Env{Key: key, Value: value})
	}
	return mergedEnvs
}

func processConfig(config *Config) {
	for i, job := range config.Jobs {
		if job.Cron == "" {
			config.Jobs[i].Cron = config.Defaults.Cron
		}
		config.Jobs[i].Envs = mergeEnvs(config.Defaults.Envs, job.Envs)
	}
}

func New(filePath string) (*Config, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s, error: %v", filePath, err)
	}

	var config Config
	if err := yaml.Unmarshal([]byte(fileData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse file: %s, error: %v", filePath, err)
	}

	processConfig(&config)
	return &config, nil
}
