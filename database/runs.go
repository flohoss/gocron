package database

func (s *Service) GetRuns(jobId uint64) []Run {
	var runs []Run
	s.orm.Where("job_id = ?", jobId).Preload("Logs").Order("start_time DESC").Find(&runs)
	return runs
}
