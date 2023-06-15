package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger(logLevel string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	zapLevel, err := zap.ParseAtomicLevel(logLevel)
	if err == nil {
		level = zapLevel
	}

	config := zap.Config{
		Level:             level,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{"pid": os.Getpid()},
	}
	return zap.Must(config.Build())
}
