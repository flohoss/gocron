package services

type Status int64

const (
	Running Status = iota + 1
	Stopped
	Finished
)
