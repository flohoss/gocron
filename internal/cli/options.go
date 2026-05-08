package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flohoss/gocron/config"
	"github.com/go-playground/validator/v10"
)

type Options struct {
	ConfigFolder string `validate:"required,dirpath"`
	ShowVersion  bool   `validate:"-"`
}

func Parse(args []string) (Options, error) {
	opts := Options{}
	flagSet := flag.NewFlagSet("gocron", flag.ContinueOnError)
	flagSet.StringVar(&opts.ConfigFolder, "config", config.GetDefaultConfigFolder(), "Path to the configuration folder")
	flagSet.BoolVar(&opts.ShowVersion, "version", false, "Print version information and exit")

	if err := flagSet.Parse(args); err != nil {
		return Options{}, err
	}

	if opts.ShowVersion {
		return opts, nil
	}

	v := validator.New()
	opts.ConfigFolder = normalizeDirPath(opts.ConfigFolder)

	if err := v.Struct(opts); err != nil {
		return Options{}, fmt.Errorf("invalid startup options: %w", err)
	}

	return opts, nil
}

func normalizeDirPath(path string) string {
	cleanPath := filepath.Clean(path)
	if cleanPath == "." {
		return "." + string(os.PathSeparator)
	}

	return cleanPath + string(os.PathSeparator)
}
