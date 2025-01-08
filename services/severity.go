package services

type Severity int

const (
	Debug Severity = iota + 1
	Info
	Warning
	Error
)
