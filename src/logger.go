package main

import (
	"gobackup/database"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogger(logLevel string) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   database.Storage + "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	zapLevel, err := zap.ParseAtomicLevel(logLevel)
	if err == nil {
		level = zapLevel
	}

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}
