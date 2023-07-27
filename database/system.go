package database

import "time"

type JobStats struct {
	TotalRuns   uint64 `json:"total_runs"`
	TotalLogs   uint64 `json:"total_logs"`
	InfoLogs    uint64 `json:"info_logs"`
	WarningLogs uint64 `json:"warning_logs"`
	ErrorLogs   uint64 `json:"error_logs"`
	ResticRuns  uint64 `json:"restic_runs"`
	CustomRuns  uint64 `json:"custom_runs"`
	PruneRuns   uint64 `json:"prune_runs"`
	CheckRuns   uint64 `json:"check_runs"`
}

func (s *Service) GetSystemLogs() []SystemLog {
	var logs []SystemLog
	s.orm.Order("ID desc").Find(&logs)
	return logs
}

func (s *Service) GetJobStats() JobStats {
	var stats JobStats
	s.orm.Model(&Run{}).Where("start_time > ?", time.Now().UnixMilli()-TimeToGoBackInMilliseconds).
		Select(`COUNT(DISTINCT runs.id) AS total_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type = ? THEN logs.run_id END) AS restic_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type = ? THEN logs.run_id END) AS custom_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type = ? THEN logs.run_id END) AS prune_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type = ? THEN logs.run_id END) AS check_runs,
				COUNT(logs.id) AS total_logs,
				SUM(CASE WHEN logs.log_severity = ? THEN 1 ELSE 0 END) AS info_logs,
				SUM(CASE WHEN logs.log_severity = ? THEN 1 ELSE 0 END) AS warning_logs,
				SUM(CASE WHEN logs.log_severity = ? THEN 1 ELSE 0 END) AS error_logs`, LogRestic, LogCustom, LogPrune, LogCheck, LogInfo, LogWarning, LogError).
		Joins("LEFT JOIN logs ON runs.id = logs.run_id").
		Scan(&stats)
	return stats
}
