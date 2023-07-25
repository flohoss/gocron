package database

type EnumLogType uint8

const (
	LogGeneral EnumLogType = iota + 1
	LogRestic
	LogCustomCommand
	LogTypePrune
	LogCheck
)

type EnumLogSeverity uint8

const (
	LogInfo EnumLogSeverity = iota + 1
	LogWarning
	LogError
)
