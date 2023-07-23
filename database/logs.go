package database

type EnumLogType uint8

const (
	LogTypeBackup EnumLogType = iota + 1
	LogTypeDocker
	LogTypePrune
	LogTypeCheck
)

type EnumLogSeverity uint8

const (
	LogSeverityInfo EnumLogSeverity = iota + 1
	LogSeverityWarning
	LogSeverityError
)
