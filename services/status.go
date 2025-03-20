package services

type Status int64

const (
	Running Status = iota + 1
	Stopped
	Finished
)

func (s Status) Int64() int64 {
	return int64(s)
}
