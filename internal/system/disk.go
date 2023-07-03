package system

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type Disk struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Unit  string  `json:"unit"`
}

func DiskUsage() Disk {
	var result = Disk{}
	d, err := disk.Usage("/")
	if err != nil {
		return result
	}
	unit, unitStr := amountString(d.Total)
	result.Unit = unitStr
	result.Total = float64(d.Total) / float64(unit)
	result.Used = float64(d.Used) / float64(unit)
	return result
}

const (
	KiB uint64 = 1024
	MiB        = KiB * 1024
	GiB        = MiB * 1024
	TiB        = GiB * 1024
	PiB        = TiB * 1024
	EiB        = PiB * 1024
)

func amountString(size uint64) (uint64, string) {
	switch {
	case size < MiB:
		return KiB, "KiB"
	case size < GiB:
		return MiB, "MiB"
	case size < TiB:
		return GiB, "GiB"
	case size < PiB:
		return TiB, "TiB"
	case size < EiB:
		return PiB, "PiB"
	default:
		return EiB, "EiB"
	}
}
