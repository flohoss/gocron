package healthcheck

type HealthCheck struct {
	Start Url `validate:"omitempty" yaml:"start"`
	End   Url `validate:"omitempty" yaml:"end"`
}

type Url struct {
	Url    string            `validate:"url" yaml:"url"`
	Params map[string]string `yaml:"params"`
	Body   string            `yaml:"body"`
}
