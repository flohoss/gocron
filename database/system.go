package database

import "time"

func (s *Service) GetSystemLogs() []SystemLog {
	var logs []SystemLog
	s.orm.Order("ID desc").Find(&logs)
	return logs
}

func (s *Service) GetJobStats() JobStats {
	var stats JobStats
	s.orm.Model(&Run{}).Where("start_time > ?", time.Now().UnixMilli()-TimeToGoBackInMilliseconds).
		Select(`COUNT(DISTINCT runs.id) AS total_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type_id = 2 THEN logs.run_id END) AS restic_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type_id = 3 THEN logs.run_id END) AS custom_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type_id = 4 THEN logs.run_id END) AS prune_runs,
				COUNT(DISTINCT CASE WHEN logs.log_type_id = 5 THEN logs.run_id END) AS check_runs,
				COUNT(logs.id) AS total_logs,
				SUM(CASE WHEN logs.log_severity_id = 1 THEN 1 ELSE 0 END) AS info_logs,
				SUM(CASE WHEN logs.log_severity_id = 2 THEN 1 ELSE 0 END) AS warning_logs,
				SUM(CASE WHEN logs.log_severity_id = 3 THEN 1 ELSE 0 END) AS error_logs`).
		Joins("LEFT JOIN logs ON runs.id = logs.run_id").
		Scan(&stats)
	return stats
}
