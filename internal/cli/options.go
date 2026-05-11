package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flohoss/gocron/config"
	"github.com/go-playground/validator/v10"
)

type Options struct {
	ConfigFile  string `validate:"required,filepath,config_file"`
	ShowVersion bool
}

func Parse(args []string) (Options, error) {
	opts := Options{}
	flagSet := flag.NewFlagSet("gocron", flag.ContinueOnError)
	flagSet.StringVar(&opts.ConfigFile, "config", config.GetDefaultConfigFile(), "Path to the configuration file")
	flagSet.BoolVar(&opts.ShowVersion, "version", false, "Print version information and exit")

	if err := flagSet.Parse(args); err != nil {
		return Options{}, err
	}

	v := validator.New()
	if err := v.RegisterValidation("config_file", validateConfigFile); err != nil {
		return Options{}, fmt.Errorf("failed to initialize startup options validator: %w", err)
	}

	if opts.ShowVersion {
		return opts, nil
	}

	opts.ConfigFile = normalizeFilePath(opts.ConfigFile)
	if err := v.Struct(opts); err != nil {
		return Options{}, fmt.Errorf("invalid startup options: %w", err)
	}

	return opts, nil
}

func normalizeFilePath(path string) string {
	return filepath.Clean(path)
}

func validateConfigFile(fl validator.FieldLevel) bool {
	path := normalizeFilePath(fl.Field().String())
	if path == "." || path == string(filepath.Separator) {
		return false
	}

	for _, part := range strings.Split(path, string(filepath.Separator)) {
		if part == ".." {
			return false
		}
	}

	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}
