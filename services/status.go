package services

type Status uint8

const (
	Running Status = iota + 1
	Stopped
	Finished
)
