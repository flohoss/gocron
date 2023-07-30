package database

import "time"

func (s *Service) GetRuns(jobId uint64) []Run {
	var runs []Run
	s.orm.Where("job_id = ?", jobId).Where("start_time > ?", time.Now().UnixMilli()-TimeToGoBackInMilliseconds).Preload("Logs").Order("start_time DESC").Find(&runs)
	return runs
}
