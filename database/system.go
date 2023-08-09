package database

type JobStats struct {
	TotalRuns   uint64 `json:"total_runs" validate:"required"`
	TotalLogs   uint64 `json:"total_logs" validate:"required"`
	InfoLogs    uint64 `json:"info_logs" validate:"required"`
	WarningLogs uint64 `json:"warning_logs" validate:"required"`
	ErrorLogs   uint64 `json:"error_logs" validate:"required"`
	GeneralRuns uint64 `json:"general_runs" validate:"required"`
	ResticRuns  uint64 `json:"restic_runs" validate:"required"`
	CustomRuns  uint64 `json:"custom_runs" validate:"required"`
	PruneRuns   uint64 `json:"prune_runs" validate:"required"`
	CheckRuns   uint64 `json:"check_runs" validate:"required"`
}

func (s *Service) GetSystemLogs() []SystemLog {
	var logs []SystemLog
	s.orm.Order("ID").Limit(5).Find(&logs)
	return logs
}

func (s *Service) GetJobStats() JobStats {
	var stats JobStats
	s.orm.Raw(`
		SELECT
			COUNT(DISTINCT run_id) AS total_runs,
			SUM(CASE WHEN subq.log_type = ? THEN 1 ELSE 0 END) AS restic_runs,
			SUM(CASE WHEN subq.log_type = ? THEN 1 ELSE 0 END) AS custom_runs,
			SUM(CASE WHEN subq.log_type = ? THEN 1 ELSE 0 END) AS prune_runs,
			SUM(CASE WHEN subq.log_type = ? THEN 1 ELSE 0 END) AS check_runs
		FROM (
			SELECT run_id, MAX(log_type) AS log_type
			FROM logs
			GROUP BY run_id
		) AS subq`, LogRestic, LogCustom, LogPrune, LogCheck).Scan(&stats)
	s.orm.Raw(`
		SELECT
			COUNT(id) AS total_logs,
			SUM(CASE WHEN log_severity = ? THEN 1 ELSE 0 END) AS info_logs,
			SUM(CASE WHEN log_severity = ? THEN 1 ELSE 0 END) AS warning_logs,
			SUM(CASE WHEN log_severity = ? THEN 1 ELSE 0 END) AS error_logs
		FROM logs`, LogInfo, LogWarning, LogError).Scan(&stats)
	return stats
}
