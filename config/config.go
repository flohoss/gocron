package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type HealthCheck struct {
	Start Url `validate:"omitempty" yaml:"start"`
	End   Url `validate:"omitempty" yaml:"end"`
}

type Url struct {
	Url    string            `validate:"url" yaml:"url"`
	Params map[string]string `yaml:"params"`
	Body   string            `yaml:"body"`
}

type Env struct {
	Key   string `validate:"required,uppercase" yaml:"key"`
	Value string `validate:"required" yaml:"value"`
}

type Command struct {
	Command string `validate:"required" yaml:"command"`
}

type Job struct {
	Name     string    `validate:"required" yaml:"name"`
	Cron     string    `validate:"omitempty,cron" yaml:"cron"`
	Envs     []Env     `validate:"omitempty,dive,required" yaml:"envs"`
	Commands []Command `validate:"required,dive,required" yaml:"commands"`
}

type Config struct {
	Defaults struct {
		Cron        string      `yaml:"cron"`
		Envs        []Env       `yaml:"envs"`
		HealthCheck HealthCheck `yaml:"healthcheck"`
	} `yaml:"defaults"`
	Jobs []Job `validate:"required,dive,required" yaml:"jobs"`
}

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}

func mergeEnvs(defaultEnvs, jobEnvs []Env) []Env {
	// Use defaultEnvs as the base
	envMap := make(map[string]int)
	mergedEnvs := make([]Env, len(defaultEnvs))
	copy(mergedEnvs, defaultEnvs)

	// Map keys in defaultEnvs to their positions
	for i, env := range defaultEnvs {
		envMap[env.Key] = i
	}

	// Process jobEnvs
	for _, env := range jobEnvs {
		if idx, exists := envMap[env.Key]; exists {
			// Override value if key already exists
			mergedEnvs[idx].Value = env.Value
		} else {
			// Append new key-value pair
			mergedEnvs = append(mergedEnvs, env)
			envMap[env.Key] = len(mergedEnvs) - 1
		}
	}

	return mergedEnvs
}

func (c *Config) processConfig() {
	for i, job := range c.Jobs {
		if job.Cron == "" {
			c.Jobs[i].Cron = c.Defaults.Cron
		}
		c.Jobs[i].Envs = mergeEnvs(c.Defaults.Envs, job.Envs)
	}
}

func readOrCreateInitFile(filePath string) ([]byte, error) {
	fileData, err := os.ReadFile(filePath)
	if err == nil {
		return fileData, nil
	}

	// File read failed; create file with default content
	defaultContent := `defaults:
  cron: '0 3 * * 0'
  envs:
    - key: SLEEP_TIME
      value: '5'
  healthcheck:
    start:
      url: http://localhost:8080
      params:
        foo: bar
      body: '{"foo": "bar"}'

jobs:
  - name: Example
    cron: '0 5 * * 0'
    commands:
      - command: ls -la
      - command: sleep ${{ SLEEP_TIME }}
      - command: echo "Done!"
      - command: sleep ${{ SLEEP_TIME }}
  - name: Example
    commands:
      - command: ls -la
      - command: sleep ${{ SLEEP_TIME }}
      - command: echo "Done!"
      - command: sleep ${{ SLEEP_TIME }}
`

	if writeErr := os.WriteFile(filePath, []byte(defaultContent), 0644); writeErr != nil {
		return nil, fmt.Errorf("failed to read file: %s, and failed to write default content: %v", filePath, writeErr)
	}

	return []byte(defaultContent), nil
}

func New(filePath string) (*Config, error) {
	fileData, err := readOrCreateInitFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s, error: %v", filePath, err)
	}

	var config Config
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse file: %s, error: %v", filePath, err)
	}

	config.processConfig()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	fmt.Printf("Config: %v\n", config.Defaults.HealthCheck.Start)
	return &config, nil
}
