package database

import "time"

func (s *Service) GetRuns() []Run {
	var runs []Run
	s.orm.Where("start_time > ?", time.Now().UnixMilli()-TimeToGoBackInMilliseconds).Preload("Logs").Order("start_time DESC").Find(&runs)
	return runs
}
