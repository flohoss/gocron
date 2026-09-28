package cli

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/validate"
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

	if opts.ShowVersion {
		return opts, nil
	}

	opts.ConfigFile = normalizeFilePath(opts.ConfigFile)
	if err := validate.Struct(opts); err != nil {
		return Options{}, fmt.Errorf("invalid startup options:\n%s", err)
	}

	return opts, nil
}

func normalizeFilePath(path string) string {
	return filepath.Clean(path)
}
