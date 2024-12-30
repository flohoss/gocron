package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Env struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Job struct {
	Name     string `yaml:"name"`
	Cron     string `yaml:"cron,omitempty"`
	Envs     []Env  `yaml:"envs"`
	Commands []struct {
		Command string `yaml:"command"`
	} `yaml:"commands"`
}

type Config struct {
	Defaults struct {
		Cron string `yaml:"cron"`
		Envs []Env  `yaml:"envs"`
	} `yaml:"defaults"`
	Jobs []Job `yaml:"jobs"`
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
