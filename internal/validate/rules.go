package validate

import (
	"net/url"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	configFileTag = "config_file"
	corsOriginTag = "cors_origin"
	envKeyTag     = "env_key"

	wildcardOrigin = "*"
)

var envKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func registerRules() {
	for _, rule := range []struct {
		tag     string
		message string
		fn      validator.Func
	}{
		{configFileTag, "{0} must be a path to a .yaml or .yml file", isConfigFile},
		{corsOriginTag, "{0} must be '*' or a scheme://host origin without a trailing slash", allowAllOriginsOrOrigin},
		{envKeyTag, "{0} must be a valid environment variable name", isValidEnvKey},
	} {
		if err := RegisterRule(rule.tag, rule.message, rule.fn); err != nil {
			panic(err)
		}
	}
}

func isConfigFile(fl validator.FieldLevel) bool {
	path := filepath.Clean(fl.Field().String())
	if path == "." || path == string(filepath.Separator) {
		return false
	}

	if slices.Contains(strings.Split(path, string(filepath.Separator)), "..") {
		return false
	}

	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

func allowAllOriginsOrOrigin(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == wildcardOrigin {
		return true
	}

	origin, err := url.Parse(value)
	if err != nil || origin.Scheme == "" || origin.Host == "" {
		return false
	}

	return origin.Path == "" && origin.RawQuery == "" && origin.Fragment == ""
}

func isValidEnvKey(fl validator.FieldLevel) bool {
	return envKeyPattern.MatchString(fl.Field().String())
}
