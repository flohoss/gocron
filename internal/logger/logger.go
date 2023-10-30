package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func NewLogger(logLevel string) *slog.Logger {
	w := os.Stderr

	lvl := new(slog.Level)
	level := slog.LevelInfo
	err := lvl.UnmarshalText([]byte(logLevel))
	if err == nil {
		level = *lvl
	}

	return slog.New(tint.NewHandler(w, &tint.Options{
		Level:      &level,
		TimeFormat: time.Kitchen,
	}))
}
